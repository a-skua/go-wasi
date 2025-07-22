package proxy

import (
	gohttp "net/http"

	"github.com/a-skua/go-wasi/http"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
)

var Handler gohttp.Handler

func init() {
	incominghandler.Exports.Handle = handle
}

type handler struct{}

func handle(in types.IncomingRequest, out types.ResponseOutparam) {
	r, err := http.ParseRequest(in)
	if err != nil {
		panic(err) // TODO: handle error properly
	}

	w := http.NewResponseWriter(out)
	defer w.Flush()

	Handler.ServeHTTP(w, r)
}
