package http

import (
	"fmt"
	"net/http"
	"net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/outgoing-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit"
)

// wasi:http/proxy
func ServeProxy(h http.Handler) error {
	incominghandler.Exports.Handle = newProxy(h)
	return nil
}

type Client http.Client

func (c *Client) Get(rawurl string) (*http.Response, error) {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	headers := types.NewFields()

	req := types.NewOutgoingRequest(headers)
	req.SetMethod(types.MethodGet())
	req.SetScheme(cm.Some(types.SchemeHTTPS())) // FIXME
	req.SetAuthority(cm.Some(url.Host))
	req.SetPathWithQuery(cm.Some(url.Path + "?" + url.RawQuery))
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
	body := &body{
		stream: wit.UnwrapResult(in.Stream()),
	}
	return &http.Response{
		StatusCode: int(res.Status()),
		Body:       body,
	}, nil
}
