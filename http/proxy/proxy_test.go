package proxy

import (
	gohttp "net/http"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/incoming-handler"
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
)

var _handler = NewHandler[
	types.IncomingRequest,
	types.ResponseOutparam,
	types.OutgoingResponse,
	types.Method,
](gohttp.HandlerFunc(
	func(w gohttp.ResponseWriter, r *gohttp.Request) {
		// TODO
	}),
	types.NewFields,
	types.NewOutgoingResponse,
	types.ResponseOutparamSet,
	types.OutgoingBodyFinish,
)

func init() {
	incominghandler.Exports.Handle = _handler.Handle
}
