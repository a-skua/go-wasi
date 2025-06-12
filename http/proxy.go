package http

import (
	"bytes"
	"fmt"
	"io"
	gohttp "net/http"
	"net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit"
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
	// FIXME
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
	con, err := wit.HandleResult(in.Consume())
	if err != nil {
		return nil, fmt.Errorf("failed to consume body: %s", err)
	}

	stream, err := wit.HandleResult(con.Stream())
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %s", err)
	}

	return &body{
		stream: stream,
	}, nil
}

func (b *body) Read(p []byte) (int, error) {
	const zero = 0
	if b.stream == nil {
		return zero, io.EOF // no body to read
	}

	list, err := wit.HandleResult(b.stream.Read(uint64(len(p))))
	if err != nil {
		return zero, fmt.Errorf("failed to read body: %s", err)
	}

	// copy
	n := int(list.Len())
	if n > len(p) {
		n = len(p)
	}
	copy(p, list.Slice())
	return n, nil
}

func (b *body) Close() error {
	b.stream.ResourceDrop()
	return nil
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

	body, err := wit.HandleResult(w.Body())
	if err != nil {
		// TODO: handle error properly
		panic(fmt.Errorf("failed to get outgoing body: %s", err))
	}
	defer types.OutgoingBodyFinish(*body, cm.None[types.Trailers]())

	output, err := wit.HandleResult(body.Write())
	if err != nil {
		// TODO: handle error properly
		panic(fmt.Errorf("failed to write body: %s", err))
	}
	defer output.ResourceDrop()

	output.Write(cm.ToList(r.body.Bytes()))
}
