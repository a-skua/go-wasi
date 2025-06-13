package main

import (
	"fmt"
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http/proxy"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
)

func init() {
	proxyHandler := proxy.NewHandler[
		types.IncomingRequest,
		types.ResponseOutparam,
		types.Method,
		types.Scheme,
		types.Headers,
		types.IncomingBody,
		types.ErrorCode,
		types.OutgoingResponse,
	](
		http.HandlerFunc(httpHandler),
		func() types.Headers {
			return types.NewFields()
		},
		func(headers types.Headers) types.OutgoingResponse {
			return types.NewOutgoingResponse(headers)
		},
		func(param types.ResponseOutparam, response types.OutgoingResponse) {
			types.ResponseOutparamSet(
				param,
				cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](response),
			)
		},
		func(body types.OutgoingBody, trailers cm.Option[types.Trailers]) {
			types.OutgoingBodyFinish(body, trailers)
		},
	)

	incominghandler.Exports.Handle = proxyHandler.Handle
}

func main() {}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	const zero = ""
	name := r.URL.Query().Get("name")
	if name == zero {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!\n", name)
	w.WriteHeader(http.StatusTeapot)
}
