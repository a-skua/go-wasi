package main

import (
	"fmt"
	"log/slog"
	gohttp "net/http"
	"os"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http"
	"github.com/a-skua/go-wasi/internal/gen/wasi/cli/run"
)

func init() {
	run.Exports.Run = runner
}

func main() {
	runner()
}

func runner() cm.BoolResult {
	var c http.Client

	res, err := c.Get("https://example.com")
	if err != nil {
		slog.Error("Failed to make GET request", "error", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	slog.Info("Response received", "status", res.Status)
	if res.StatusCode != gohttp.StatusOK {
		slog.Error("Unexpected status code", "status_code", res.StatusCode)
		os.Exit(1)
	}

	body := make([]byte, 1024)
	_, err = res.Body.Read(body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		os.Exit(1)
	}

	fmt.Println(string(body))

	return cm.ResultOK
}
