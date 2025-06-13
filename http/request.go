package http

import (
	"fmt"
	goio "io"
	"net/http"
	"net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit"
)

func ParseRequest(in types.IncomingRequest) (*http.Request, error) {
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

func parseUrl(in types.IncomingRequest) (*url.URL, error) {
	schemeOpt := in.Scheme()
	if schemeOpt.None() {
		return nil, fmt.Errorf("scheme is required")
	}
	scheme := schemeOpt.Value()

	authority, ok := wit.HandleOption(in.Authority())
	if !ok {
		return nil, fmt.Errorf("authority is required")
	}

	path := wit.UnwrapOptionOr(in.PathWithQuery(), "/")

	rawURL := fmt.Sprintf("%s://%s%s",
		scheme.String(),
		authority,
		path,
	)

	return url.ParseRequestURI(rawURL)
}

type body struct {
	stream types.InputStream
}

func parseBody(in types.IncomingRequest) (*body, error) {
	con, err := wit.HandleResult(in.Consume())
	if err != nil {
		return nil, fmt.Errorf("failed to consume body: %s", err)
	}

	stream, err := wit.HandleResult((*con).Stream())
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %s", err)
	}

	return &body{
		stream: *stream,
	}, nil
}

func (b *body) Read(p []byte) (int, error) {
	const zero = 0
	if b == nil {
		return zero, goio.EOF
	}

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

func (b *body) Close() error {
	b.stream.ResourceDrop()
	return nil
}

func parseHeaders(in types.IncomingRequest) http.Header {
	headers := http.Header{}

	entries := in.Headers().Entries()
	for _, entry := range entries.Slice() {
		k := string(entry.F0)
		v := string(cm.List[uint8](entry.F1).Slice())
		headers[k] = append(headers[k], v)
	}
	return headers
}

type Header struct {
	http.Header
	Status int
}

func NewHeader() Header {
	return Header{
		Header: make(http.Header),
		Status: 200,
	}
}
