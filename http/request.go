package http

import (
	"io"
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

func newRequest(r *http.Request) (zero types.OutgoingRequest, _ error) {
	out := types.NewOutgoingRequest(newHeader(r.Header).headers())
	out.SetMethod(types.Method(newMethod(r.Method)))

	err := url.SetOutgoingRequestURL(out, r.URL)
	if err != nil {
		return zero, err
	}

	if r.Body == nil {
		return out, nil
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return zero, err
	}

	body := result.Unwrap(out.Body())
	defer types.OutgoingBodyFinish(body, cm.None[types.Trailers]())

	stream := result.Unwrap(body.Write())
	stream.Write(cm.ToList(buf))
	defer stream.ResourceDrop()

	return out, err
}
