package http

import (
	"net/http"
	gourl "net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http/internal/url"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/outgoing-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/future"
	"github.com/a-skua/go-wasi/internal/wit/option"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

type Client http.Client

func (c *Client) Get(rawurl string) (*http.Response, error) {
	h := newHeader()

	out := types.NewOutgoingRequest(h.headers())
	out.SetMethod(types.MethodGet())

	u, err := gourl.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}
	url.SetOutgoingRequestURL(out, u)

	f, err := result.Handle(outgoinghandler.Handle(out, cm.None[types.RequestOptions]()))
	if err != nil {
		return nil, err
	}
	defer f.ResourceDrop()

	future.Wait(f)

	r, err := result.Handle(result.Unwrap(option.Unwrap(f.Get())))
	if err != nil {
		return nil, err
	}
	defer r.ResourceDrop()

	return parseResponse(r)
}
