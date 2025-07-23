package main

import (
	"fmt"
	"net/http"

	"github.com/a-skua/go-wasi/http/proxy"
)

func init() {
	proxy.Handler = http.HandlerFunc(httpHandler)
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
