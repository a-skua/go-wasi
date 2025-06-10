package http

import (
	"net/http"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
)

// wasi:http/proxy
func ServeProxy(h http.Handler) error {
	incominghandler.Exports.Handle = newProxy(h)
	return nil
}
