package main

import (
	"fmt"
	"net/http"

	proxy "github.com/a-skua/go-wasi/http"
)

func init() {
	proxy.ServeProxy(http.HandlerFunc(handler))
}

func main() {}

func handler(w http.ResponseWriter, r *http.Request) {
	const zero = ""
	name := r.URL.Query().Get("name")
	if name == zero {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!\n", name)
	w.WriteHeader(http.StatusTeapot)
}
