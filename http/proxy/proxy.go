package proxy

import (
	"bytes"
	gohttp "net/http"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/internal/wit/result"
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

	w := newResponse()
	defer w.flush(out)

	Handler.ServeHTTP(w, r)
}

type response struct {
	header http.Header
	body   bytes.Buffer
}

func newResponse() *response {
	return &response{
		header: http.NewHeader(),
	}
}

func (r *response) Header() gohttp.Header {
	return r.header.Header
}

func (r *response) Write(b []byte) (int, error) {
	r.body.Write(b)
	return len(b), nil
}

func (r *response) WriteHeader(statusCode int) {
	r.header.Status = statusCode
}

func newHeaders() types.Headers {
	return types.NewFields()
}

func (r *response) flush(out types.ResponseOutparam) {
	headers := newHeaders()
	for k, vs := range r.header.Header {
		if vs == nil {
			continue
		}
		for _, v := range vs {
			headers.Append(types.FieldName(k), types.FieldValue(cm.ToList([]uint8(v))))
		}
	}

	w := types.NewOutgoingResponse(headers)
	w.SetStatusCode(types.StatusCode(r.header.Status))

	defer types.ResponseOutparamSet(out, cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](w))

	body := result.Unwrap(w.Body())
	defer types.OutgoingBodyFinish(body, cm.None[types.Trailers]())

	output := result.Unwrap(body.Write())
	defer (output).ResourceDrop()

	output.Write(cm.ToList(r.body.Bytes()))
}
