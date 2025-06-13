package proxy

import (
	gohttp "net/http"
)

func init() {
	Handler = NewHandler(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		// TODO
	}))
}
