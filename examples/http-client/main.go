package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	gohttp "net/http"
	"os"
	"strings"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http"
	"github.com/a-skua/go-wasi/internal/gen/wasi/cli/run"
)

var (
	method = flag.String("x", "GET", "GET or POST")
	url    = flag.String("url", "https://example.com", "URL to fetch")
)

func init() {
	flag.Parse()
	run.Exports.Run = Run
}

func main() {}

func Run() cm.BoolResult {
	var (
		c   http.Client
		r   *gohttp.Response
		err error
	)
	if *method == "GET" {
		r, err = c.Get(*url)
	} else {
		r, err = c.Post(*url, "text/plan", strings.NewReader("hello, world"))
	}
	if err != nil {
		slog.Error("Failed to make GET request", "error", err)
		os.Exit(1)
	}
	defer r.Body.Close()

	if r.StatusCode != gohttp.StatusOK {
		slog.Error("Unexpected status code", "status_code", r.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		os.Exit(1)
	}

	fmt.Println(string(body))

	return cm.ResultOK
}
