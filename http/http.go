package http

import (
	gohttp "net/http"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
)

// wasi:http/proxy
func ServeProxy(h gohttp.Handler) error {
	incominghandler.Exports.Handle = newProxy(h)
	return nil
}
