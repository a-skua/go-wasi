package http

import (
	"io"
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/outgoing-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/future"
	"github.com/a-skua/go-wasi/internal/wit/option"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

type Client http.Client

func (c *Client) Get(rawurl string) (zero *http.Response, _ error) {
	r, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return zero, err
	}

	return c.Do(r)
}

func (c *Client) Post(rawurl, contentType string, body io.Reader) (zero *http.Response, _ error) {
	r, err := http.NewRequest(http.MethodPost, rawurl, body)
	if err != nil {
		return zero, err
	}
	r.Header.Set("Content-Type", contentType)

	return c.Do(r)
}

func (c *Client) Do(r *http.Request) (zero *http.Response, _ error) {
	out, err := newRequest(r)
	if err != nil {
		return zero, err
	}

	o := types.NewRequestOptions()

	f, err := result.Handle(outgoinghandler.Handle(out, cm.Some(o)))
	if err != nil {
		return zero, err
	}
	future.Wait(f)

	in, err := result.Handle(result.Unwrap(option.Unwrap(f.Get())))
	if err != nil {
		return zero, err
	}

	return parseResponse(in)
}
