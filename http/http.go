package http

import (
	"io"
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

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

type body struct {
	stream types.InputStream
}

func parseBody(in types.IncomingBody) (*body, error) {
	stream, err := result.Handle(in.Stream())
	if err != nil {
		return nil, err
	}

	return &body{
		stream: stream,
	}, nil
}

func (b *body) Read(p []byte) (zero int, _ error) {
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

func (b *body) Close() error {
	b.stream.ResourceDrop()
	return nil
}
