package http

import (
	"bytes"
	"fmt"
	"io"
	gohttp "net/http"
	"net/url"

	"github.com/a-skua/go-wasi/internal/wasi/http/types"
	"go.bytecodealliance.org/cm"
)

type proxyHandler func(request types.IncomingRequest, responseOut types.ResponseOutparam)

func newProxy(handler gohttp.Handler) proxyHandler {
	return func(in types.IncomingRequest, out types.ResponseOutparam) {
		r, err := parseRequest(in)
		if err != nil {
			panic(err) // TODO: handle error properly
		}

		w := newResponse()
		defer w.flush(out)

		handler.ServeHTTP(w, r)
	}
}

func parseUrl(in types.IncomingRequest) (*url.URL, error) {
	rawURL := fmt.Sprintf("%s://%s%s",
		in.Scheme().Value(),
		in.Authority().Value(),
		in.PathWithQuery().Value(),
	)

	return url.ParseRequestURI(rawURL)
}

type body struct {
	stream *types.InputStream
}

func parseBody(in types.IncomingRequest) (*body, error) {
	con := in.Consume()
	if con.IsErr() {
		return nil, fmt.Errorf("failed to consume body: %s", con.Err())
	}

	stream := con.OK().Stream()
	if stream.IsErr() {
		return nil, fmt.Errorf("failed to get stream: %s", stream.Err())
	}

	return &body{
		stream: stream.OK(),
	}, nil
}

func (b *body) Read(p []byte) (int, error) {
	if b.stream == nil {
		return 0, io.EOF // no body to read
	}

	result := b.stream.Read(uint64(len(p)))
	if result.IsErr() {
		return 0, fmt.Errorf("failed to read body: %s", result.Err())
	}

	// copy
	n := int(result.OK().Len())
	if n > len(p) {
		n = len(p)
	}
	copy(p, result.OK().Slice())
	return n, nil
}

func parseHeaders(in types.IncomingRequest) gohttp.Header {
	headers := gohttp.Header{}

	for _, t := range in.Headers().Entries().Slice() {
		k := string(t.F0)
		v := string(t.F1.Slice())
		headers[k] = append(headers[k], v)
	}
	return headers
}

func parseRequest(in types.IncomingRequest) (*gohttp.Request, error) {
	method := in.Method()

	url, err := parseUrl(in)
	if err != nil {
		return nil, err
	}

	body, err := parseBody(in)
	if err != nil {
		return nil, err
	}

	r, err := gohttp.NewRequest(method.String(), url.String(), body)
	if err != nil {
		return nil, err
	}

	r.Header = parseHeaders(in)

	return r, nil
}

type header struct {
	gohttp.Header
	status int
}

func newHeader() header {
	return header{
		Header: make(gohttp.Header),
		status: 200,
	}
}

func (h header) headers() types.Headers {
	headers := types.NewFields()
	for k, vs := range h.Header {
		if vs == nil {
			continue // skip nil values
		}
		for _, v := range vs {
			headers.Append(types.FieldKey(k), types.FieldValue(cm.ToList([]byte(v))))
		}
	}
	return headers
}

type response struct {
	status int
	header header
	body   bytes.Buffer
}

func newResponse() *response {
	return &response{
		header: newHeader(),
	}
}

func (r *response) Header() gohttp.Header {
	return r.header.Header
}

func (r *response) Write(b []byte) (int, error) {
	r.body.Write(b)
	return len(b), nil
}

func (r *response) WriteHeader(statusCode int) {
	r.header.status = statusCode
}

func (r *response) flush(out types.ResponseOutparam) {
	w := types.NewOutgoingResponse(r.header.headers())
	w.SetStatusCode(types.StatusCode(r.header.status))
	defer types.ResponseOutparamSet(
		out,
		cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](w),
	)

	if b := w.Body(); b.IsOK() {
		body := *b.OK()
		defer types.OutgoingBodyFinish(body, cm.None[types.Trailers]())

		if w := body.Write(); w.IsOK() {
			output := w.OK()
			defer output.ResourceDrop()

			output.Write(cm.ToList(r.body.Bytes()))
		}
	}
}
