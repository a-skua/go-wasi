package http

import (
	"io"
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http/internal/url"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

func ParseRequest(in types.IncomingRequest) (*http.Request, error) {
	method := in.Method()

	url, err := url.ParseIncomingRequest(in)
	if err != nil {
		return nil, err
	}

	body, err := parseRequestBody(in)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method.String(), url.String(), body)
	if err != nil {
		return nil, err
	}

	r.Header = parseRequestHeaders(in)

	return r, nil
}

type requestBody struct {
	stream types.InputStream
}

func parseRequestBody(in types.IncomingRequest) (*requestBody, error) {
	con, err := result.Handle(in.Consume())
	if err != nil {
		return nil, err
	}

	stream, err := result.Handle(con.Stream())
	if err != nil {
		return nil, err
	}

	return &requestBody{
		stream: stream,
	}, nil
}

func (b *requestBody) Read(p []byte) (zero int, _ error) {
	if b == nil {
		return zero, io.EOF
	}

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

func (b *requestBody) Close() error {
	b.stream.ResourceDrop()
	return nil
}

func parseRequestHeaders(in types.IncomingRequest) http.Header {
	headers := http.Header{}

	entries := in.Headers().Entries()
	for _, entry := range entries.Slice() {
		k := string(entry.F0)
		v := string(cm.List[uint8](entry.F1).Slice())
		headers[k] = append(headers[k], v)
	}
	return headers
}
