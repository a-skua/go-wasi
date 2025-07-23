package http

import (
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http/internal/url"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

func ParseRequest(r types.IncomingRequest) (*http.Request, error) {
	method := r.Method()

	url, err := url.ParseIncomingRequest(r)
	if err != nil {
		return nil, err
	}

	in, err := result.Handle(r.Consume())
	if err != nil {
		return nil, err
	}

	body, err := parseBody(in)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method.String(), url.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header = parseHeaders(r.Headers())

	return request, nil
}

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
