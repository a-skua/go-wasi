package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/outgoing-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit"
)

type Client http.Client

type clientBody struct {
	stream types.InputStream
}

func (b *clientBody) Read(p []byte) (int, error) {
	const zero = 0
	list, err := wit.HandleResult(b.stream.Read(uint64(len(p))))
	if err != nil {
		return zero, fmt.Errorf("failed to read body: %s", err)
	}

	n := int(list.Len())
	if n > len(p) {
		n = len(p)
	}
	copy(p, list.Slice())
	return n, nil
}

func (b *clientBody) Close() error {
	b.stream.ResourceDrop()
	return nil
}

func (c *Client) Get(rawurl string) (*http.Response, error) {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	headers := types.NewFields()

	req := types.NewOutgoingRequest(headers)
	req.SetMethod(types.MethodGet())

	switch url.Scheme {
	case "http":
		req.SetScheme(cm.Some(types.SchemeHTTP()))
	case "https":
		req.SetScheme(cm.Some(types.SchemeHTTPS()))
	default:
		req.SetScheme(cm.Some(types.SchemeHTTPS()))
	}

	req.SetAuthority(cm.Some(url.Host))
	pathWithQuery := url.Path
	if url.RawQuery != "" {
		pathWithQuery += "?" + url.RawQuery
	}
	req.SetPathWithQuery(cm.Some(pathWithQuery))
	future, errcode := wit.HandleResult(outgoinghandler.Handle(req, cm.None[types.RequestOptions]()))
	if errcode != nil {
		return nil, fmt.Errorf("failed to handle outgoing request: %v", errcode)
	}
	defer future.ResourceDrop()

	poll := future.Subscribe()
	defer poll.ResourceDrop()
	poll.Block()

	wrap := wit.UnwrapResult(future.Get().Value())
	res, errcode := wit.HandleResult(*wrap)
	if errcode != nil {
		return nil, fmt.Errorf("failed to get future response: %v", errcode)
	}
	defer res.ResourceDrop()

	in := wit.UnwrapResult(res.Consume())
	stream := wit.UnwrapResult(in.Stream())
	clientBody := &clientBody{
		stream: *stream,
	}
	return &http.Response{
		StatusCode: int(res.Status()),
		Body:       clientBody,
	}, nil
}
func parseRequest(in types.IncomingRequest) (*http.Request, error) {
	method := in.Method()

	url, err := parseUrl(in)
	if err != nil {
		return nil, err
	}

	body, err := parseBody(in)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method.String(), url.String(), body)
	if err != nil {
		return nil, err
	}

	r.Header = parseHeaders(in)

	return r, nil
}

type proxyHeader struct {
	http.Header
	status int
}

func newHeader() proxyHeader {
	return proxyHeader{
		Header: make(http.Header),
		status: 200,
	}
}

func (h proxyHeader) headers() types.Headers {
	headers := types.NewFields()
	for k, vs := range h.Header {
		if vs == nil {
			continue
		}
		for _, v := range vs {
			headers.Append(types.FieldKey(k), types.FieldValue(cm.ToList([]byte(v))))
		}
	}
	return headers
}

type proxyResponse struct {
	status int
	header proxyHeader
	body   bytes.Buffer
}

func newResponse() *proxyResponse {
	return &proxyResponse{
		header: newHeader(),
	}
}

func (r *proxyResponse) Header() http.Header {
	return r.header.Header
}

func (r *proxyResponse) Write(b []byte) (int, error) {
	r.body.Write(b)
	return len(b), nil
}

func (r *proxyResponse) WriteHeader(statusCode int) {
	r.header.status = statusCode
}

func (r *proxyResponse) flush(out types.ResponseOutparam) {
	w := types.NewOutgoingResponse(r.header.headers())
	w.SetStatusCode(types.StatusCode(r.header.status))
	defer types.ResponseOutparamSet(
		out,
		cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](w),
	)

	body, err := wit.HandleResult(w.Body())
	if err != nil {
		panic(fmt.Errorf("failed to get outgoing body: %s", err))
	}
	defer types.OutgoingBodyFinish(*body, cm.None[types.Trailers]())

	output, err := wit.HandleResult(body.Write())
	if err != nil {
		panic(fmt.Errorf("failed to write body: %s", err))
	}
	defer output.ResourceDrop()

	output.Write(cm.ToList(r.body.Bytes()))
}
