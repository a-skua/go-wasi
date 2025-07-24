package http

import (
	"bytes"
	"fmt"
	"net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/result"
)

type ResponseWriter interface {
	http.ResponseWriter
	Flush() error
}

func NewResponseWriter(out types.ResponseOutparam) ResponseWriter {
	return &response{
		out:    out,
		header: newHeader(http.Header{}),
	}
}

type response struct {
	out    types.ResponseOutparam
	status int
	header header
	body   bytes.Buffer
}

func (r *response) Header() http.Header {
	return http.Header(r.header)
}

func (r *response) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *response) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *response) Flush() error {
	w := types.NewOutgoingResponse(r.header.headers())

	ok := result.HandleBool(w.SetStatusCode(types.StatusCode(r.status)))
	if !ok {
		return fmt.Errorf("failed to set status code %d", r.status)
	}

	defer types.ResponseOutparamSet(
		r.out,
		cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](w),
	)

	body := result.Unwrap(w.Body())
	defer types.OutgoingBodyFinish(body, cm.None[types.Trailers]())

	output := result.Unwrap(body.Write())
	defer (output).ResourceDrop()

	_, err := result.Handle(output.Write(cm.ToList(r.body.Bytes())))
	return err
}

func parseResponse(r types.IncomingResponse) (*http.Response, error) {
	status := r.Status()
	if status < 100 || status > 599 {
		return nil, fmt.Errorf("invalid status code: %d", status)
	}

	in, err := result.Handle(r.Consume())
	if err != nil {
		return nil, err
	}

	body, err := parseBody(in)
	if err != nil {
		return nil, err
	}

	response := &http.Response{
		StatusCode: int(status),
		Body:       body,
	}

	response.Header = parseHeaders(r.Headers())

	return response, nil
}
