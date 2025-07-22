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
		header: newHeader(),
	}
}

type response struct {
	out    types.ResponseOutparam
	header header
	body   bytes.Buffer
}

func (r *response) Header() http.Header {
	return r.header.Header
}

func (r *response) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *response) WriteHeader(statusCode int) {
	r.header.Status = statusCode
}

func (r *response) Flush() error {
	w := types.NewOutgoingResponse(r.header.headers())

	ok := result.HandleBool(w.SetStatusCode(types.StatusCode(r.header.Status)))
	if !ok {
		return fmt.Errorf("failed to set status code %d", r.header.Status)
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
