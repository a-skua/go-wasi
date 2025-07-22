package http

import (
	"fmt"
	"net/http"
	gourl "net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http/internal/url"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/outgoing-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/option"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

type Client http.Client

type clientBody struct {
	stream types.InputStream
}

func (b *clientBody) Read(p []byte) (zero int, _ error) {
	list, err := result.Handle(b.stream.Read(uint64(len(p))))
	if err != nil {
		return zero, err
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
	h := newHeader()

	out := types.NewOutgoingRequest(h.headers())
	out.SetMethod(types.MethodGet())

	u, err := gourl.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}
	url.SetOutgoingRequestURL(out, u)

	future, err := result.Handle(outgoinghandler.Handle(out, cm.None[types.RequestOptions]()))
	if err != nil {
		return nil, err
	}
	defer future.ResourceDrop()

	poll := future.Subscribe()
	defer poll.ResourceDrop()
	poll.Block()

	wrap := result.Unwrap(option.Unwrap(future.Get()))
	res, errcode := result.Handle(wrap)
	if errcode != nil {
		return nil, fmt.Errorf("failed to get future response: %v", errcode)
	}
	defer res.ResourceDrop()

	in := result.Unwrap(res.Consume())
	stream := result.Unwrap(in.Stream())
	clientBody := &clientBody{
		stream: stream,
	}
	return &http.Response{
		StatusCode: int(res.Status()),
		Body:       clientBody,
	}, nil
}

type header struct {
	http.Header
	status int
}

func newHeader() header {
	return header{
		Header: make(http.Header),
		status: 200,
	}
}

func (h header) headers() types.Headers {
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

func (h header) statusCode() types.StatusCode {
	return types.StatusCode(h.status)
}
