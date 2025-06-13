package http

import (
	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
	"github.com/a-skua/go-wasi/io"
)

type Scheme interface {
	HTTP() bool
	HTTPS() bool
	Other() *string
	String() string
}

type Method interface {
	Get() bool
	Head() bool
	Post() bool
	Put() bool
	Delete() bool
	Connect() bool
	Options() bool
	Trace() bool
	Patch() bool
	Other() *string
	String() string
}

type FieldKey = string
type FieldName = FieldKey
type FieldValue = cm.AnyList[uint8]

type Fields[N ~FieldName, V FieldValue, E ~HeaderError] interface {
	ResourceDrop()
	Append(name N, value V) cm.Result[E, struct{}, E]
	Delete(name N) cm.Result[E, struct{}, E]
	Entries() cm.List[cm.Tuple[N, V]]
	Get(name N) cm.List[V]
	Has(name N) bool
	Set(name N, value cm.List[V]) cm.Result[E, struct{}, E]
}

type HeaderError = uint8

// TODO
// type Headers[N FieldName, V FieldValue, E HeaderError] = Fields[N, V, E]

// TODO
// type Trailers[N FieldName, V FieldValue, E HeaderError] = Fields[N, V, E]

type IncomingBody[S io.InputStream[E, P],
	E io.StreamError, P io.Pollable,
] interface {
	ResourceDrop()
	Stream() cm.Result[S, S, struct{}]
}
type IncomingRequest[
	S Scheme, M Method, H Fields[N, V, HE], B IncomingBody[I, SE, P],
	N FieldName, V FieldValue, HE HeaderError,
	I io.InputStream[SE, P],
	SE io.StreamError, P io.Pollable,
] interface {
	ResourceDrop()
	Authority() cm.Option[string]
	Consume() cm.Result[B, B, struct{}]
	Headers() H
	Method() M
	PathWithQuery() cm.Option[string]
	Scheme() cm.Option[S]
}

type ErrorCode interface{}

type ErrorCodeShape interface{}

type ResponseOutparam[
	Headers Fields[N, V, E], Err ErrorCode,
	N FieldName, V FieldValue, E HeaderError,
] interface {
	ResourceDrop()
	SendInformational(status uint16, headers Headers) cm.Result[Err, struct{}, Err]
}

type OutgoingResponse[
	H Fields[N, V, HE],
	B OutgoingBody[O, I, SE, P],
	N FieldName, V FieldValue, HE HeaderError,
	O io.OutputStream[I, SE, P], I io.InputStream[SE, P], SE io.StreamError, P io.Pollable,
] interface {
	ResourceDrop()
	Body() cm.Result[B, B, struct{}]
	Headers() H
	SetStatusCode(statusCode StatusCode) cm.BoolResult
	StatusCode() StatusCode
}

type StatusCode = types.StatusCode

type OutgoingBody[
	O io.OutputStream[I, E, P],
	I io.InputStream[E, P], E io.StreamError, P io.Pollable,
] interface {
	ResourceDrop()
	Write() cm.Result[O, O, struct{}]
}
