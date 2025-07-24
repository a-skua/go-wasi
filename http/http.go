package http

import (
	"io"
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

func newMethod(m string) types.Method {
	switch m {
	case http.MethodGet:
		return types.MethodGet()
	case http.MethodPost:
		return types.MethodPost()
	case http.MethodPut:
		return types.MethodPut()
	case http.MethodDelete:
		return types.MethodDelete()
	case http.MethodHead:
		return types.MethodHead()
	case http.MethodPatch:
		return types.MethodPatch()
	case http.MethodOptions:
		return types.MethodOptions()
	case http.MethodTrace:
		return types.MethodTrace()
	case http.MethodConnect:
		return types.MethodConnect()
	default:
		return types.MethodOther(m)
	}
}

type header http.Header

func parseHeaders(h types.Headers) http.Header {
	headers := http.Header{}

	entries := h.Entries()
	for _, entry := range entries.Slice() {
		k := string(entry.F0)
		v := string(cm.List[uint8](entry.F1).Slice())
		headers[k] = append(headers[k], v)
	}
	return headers
}

func newHeader(h http.Header) header {
	return header(h)
}

func (h header) headers() types.Headers {
	headers := types.NewFields()
	for k, vs := range h {
		if vs == nil {
			continue
		}
		for _, v := range vs {
			headers.Append(types.FieldKey(k), types.FieldValue(cm.ToList([]byte(v))))
		}
	}
	return headers
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
