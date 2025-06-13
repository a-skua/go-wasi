package proxy

import (
	"bytes"
	"fmt"
	goio "io"
	gohttp "net/http"
	"net/url"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/http"
	"github.com/a-skua/go-wasi/internal/wit"
	"github.com/a-skua/go-wasi/io"
)

type Handler[
	Request http.IncomingRequest[
		Scheme, Method, Headers, Body,
		Fields, FieldName, FieldValue, HeaderError,
		InputStream,
		StreamError, Pollable,
	],
	Response http.ResponseOutparam[Headers, ErrorCode, Fields, FieldName, FieldValue, HeaderError],
	Method http.Method,
	Scheme http.Scheme,
	Headers http.Fields[Fields, FieldName, FieldValue, HeaderError],
	Body http.IncomingBody[InputStream, StreamError, Pollable],
	ErrorCode http.ErrorCode,
	Fields any, FieldName http.FieldName, FieldValue http.FieldValue, HeaderError http.HeaderError,
	InputStream io.InputStream[StreamError, Pollable],
	StreamError io.StreamError, Pollable io.Pollable,
] interface {
	Handle(Request, Response)
}

func NewHandler[
	Request http.IncomingRequest[
		Scheme, Method, Headers, Body,
		Fields, FieldName, FieldValue, HeaderError,
		InputStream,
		StreamError, Pollable,
	],
	Response http.ResponseOutparam[Headers, ErrorCode, FieldName, FieldValue, HeaderError],
	OutgoingResponse http.OutgoingResponse[
		Headers, OutgoingBody,
		Fields, FieldName, FieldValue, HeaderError,
		OutputStream, InputStream, StreamError, Pollable,
	],
	Method http.Method,
	Scheme http.Scheme,
	Headers http.Fields[Fields, FieldName, FieldValue, HeaderError],
	Body http.IncomingBody[InputStream, StreamError, Pollable],
	ErrorCode http.ErrorCode, ErrorCodeShape http.ErrorCodeShape,
	Fields any, FieldName http.FieldName, FieldValue http.FieldValue, HeaderError http.HeaderError,
	InputStream io.InputStream[StreamError, Pollable],
	OutgoingBody http.OutgoingBody[OutputStream, InputStream, StreamError, Pollable],
	OutputStream io.OutputStream[InputStream, StreamError, Pollable],
	Trailers http.Fields[Fields, FieldName, FieldValue, HeaderError],
	StreamError io.StreamError, Pollable io.Pollable,
](
	h gohttp.Handler,
	newHeaders func() Headers,
	newOutgoingResponse func(Headers) OutgoingResponse,
	responseOutparamSet func(Response, cm.Result[ErrorCodeShape, OutgoingResponse, ErrorCode]),
	outgoingBodyFinish func(OutgoingBody, cm.Option[Trailers]) cm.Result[ErrorCode, struct{}, ErrorCode],
) Handler[
	Request, Response, Method, Scheme, Headers, Body, ErrorCode,
	FieldName, FieldValue, HeaderError,
	InputStream,
	StreamError, Pollable,
] {
	return &handler[
		Request, Response,
		OutgoingResponse,
		Method, Scheme, Headers, Body, Trailers,
		ErrorCode, ErrorCodeShape,
		FieldName, FieldValue, HeaderError,
		OutgoingBody, OutputStream, InputStream, StreamError, Pollable,
	]{
		handler:             h,
		newHeaders:          newHeaders,
		newOutgoingResponse: newOutgoingResponse,
		responseOutparamSet: responseOutparamSet,
		outgoingBodyFinish:  outgoingBodyFinish,
	}
}

type handler[
	Request http.IncomingRequest[
		Scheme, Method, Headers, Body,
		FieldName, FieldValue, HeaderError,
		InputStream,
		StreamError, Pollable,
	],
	Response http.ResponseOutparam[
		Headers, ErrorCode,
		FieldName, FieldValue, HeaderError,
	],
	OutgoingResponse http.OutgoingResponse[
		Headers, OutgoingBody,
		FieldName, FieldValue, HeaderError,
		OutputStream, InputStream, StreamError, Pollable,
	],
	Method http.Method,
	Scheme http.Scheme,
	Headers http.Fields[FieldName, FieldValue, HeaderError],
	Body http.IncomingBody[InputStream, StreamError, Pollable],
	Trailers http.Fields[FieldName, FieldValue, HeaderError],
	ErrorCode http.ErrorCode, ErrorCodeShape http.ErrorCodeShape,
	FieldName http.FieldName,
	FieldValue http.FieldValue,
	HeaderError http.HeaderError,
	OutgoingBody http.OutgoingBody[OutputStream, InputStream, StreamError, Pollable],
	OutputStream io.OutputStream[InputStream, StreamError, Pollable],
	InputStream io.InputStream[StreamError, Pollable],
	StreamError io.StreamError, Pollable io.Pollable,
] struct {
	handler             gohttp.Handler
	newHeaders          func() Headers
	newOutgoingResponse func(Headers) OutgoingResponse
	responseOutparamSet func(Response, cm.Result[ErrorCodeShape, OutgoingResponse, ErrorCode])
	outgoingBodyFinish  func(OutgoingBody, cm.Option[Trailers]) cm.Result[ErrorCode, struct{}, ErrorCode]
}

func (h *handler[
	Request, Response, _, _, _, Headers, _, Trailers,
	_, _, _, _, _, _, _, _, _, _,
]) Handle(in Request, out Response) {
	r, err := parseRequest(in)
	if err != nil {
		panic(err) // TODO: handle error properly
	}

	w := newResponse[Headers, Trailers, Response](h.newHeaders, h.newOutgoingResponse, h.responseOutparamSet, h.outgoingBodyFinish)
	defer w.flush(out)

	h.handler.ServeHTTP(w, r)
}

func parseUrl[
	S http.Scheme, M http.Method, H http.Fields[N, V, HE], B http.IncomingBody[I, SE, P],
	N http.FieldName, V http.FieldValue, HE http.HeaderError,
	I io.InputStream[SE, P],
	SE io.StreamError, P io.Pollable,
](in http.IncomingRequest[S, M, H, B, N, V, HE, I, SE, P]) (*url.URL, error) {
	schemeOpt := in.Scheme()
	if schemeOpt.None() {
		return nil, fmt.Errorf("scheme is required")
	}
	scheme := schemeOpt.Value()

	authority, ok := wit.HandleOption(in.Authority())
	if !ok {
		return nil, fmt.Errorf("authority is required")
	}

	path := wit.UnwrapOptionOr(in.PathWithQuery(), "/")

	rawURL := fmt.Sprintf("%s://%s%s",
		scheme.String(),
		authority,
		path,
	)

	return url.ParseRequestURI(rawURL)
}

type body[
	InputStream io.InputStream[StreamError, Pollable],
	StreamError io.StreamError, Pollable io.Pollable,
] struct {
	stream InputStream
}

func parseBody[
	S http.Scheme, M http.Method, H http.Fields[N, V, HE], B http.IncomingBody[I, SE, P],
	N http.FieldName, V http.FieldValue, HE http.HeaderError,
	I io.InputStream[SE, P],
	SE io.StreamError, P io.Pollable,

](in http.IncomingRequest[S, M, H, B, N, V, HE, I, SE, P]) (*body[I, SE, P], error) {
	con, err := wit.HandleResult(in.Consume())
	if err != nil {
		return nil, fmt.Errorf("failed to consume body: %s", err)
	}

	stream, err := wit.HandleResult((*con).Stream())
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %s", err)
	}

	return &body[I, SE, P]{
		stream: *stream,
	}, nil
}

func (b *body[
	InputStream,
	StreamError, Pollable,
]) Read(p []byte) (int, error) {
	const zero = 0
	if b == nil {
		return zero, goio.EOF
	}

	list, err := wit.HandleResult(b.stream.Read(uint64(len(p))))
	if err != nil {
		return zero, fmt.Errorf("failed to read body: %s", err)
	}

	n := int(list.Len())
	if n > len(p) {
		n = len(p)
	}
	copy(p, list.Slice())
	return n, nil
}

func (b *body[
	InputStream,
	StreamError, Pollable,
]) Close() error {
	b.stream.ResourceDrop()
	return nil
}

func parseHeaders[
	S http.Scheme, M http.Method, H http.Fields[N, V, HE], B http.IncomingBody[I, SE, P],
	N http.FieldName, V http.FieldValue, HE http.HeaderError,
	I io.InputStream[SE, P],
	SE io.StreamError, P io.Pollable,
](in http.IncomingRequest[S, M, H, B, N, V, HE, I, SE, P]) gohttp.Header {
	headers := gohttp.Header{}

	entries := in.Headers().Entries()
	for _, entry := range entries.Slice() {
		k := string(entry.F0)
		v := string(cm.List[uint8](entry.F1).Slice())
		headers[k] = append(headers[k], v)
	}
	return headers
}

func parseRequest[
	S http.Scheme, M http.Method, H http.Fields[N, V, HE], B http.IncomingBody[I, SE, P],
	N http.FieldName, V http.FieldValue, HE http.HeaderError,
	I io.InputStream[SE, P],
	SE io.StreamError, P io.Pollable,
](in http.IncomingRequest[S, M, H, B, N, V, HE, I, SE, P]) (*gohttp.Request, error) {
	method := in.Method()

	url, err := parseUrl(in)
	if err != nil {
		return nil, err
	}

	body, err := parseBody(in)
	if err != nil {
		return nil, err
	}

	r, err := gohttp.NewRequest(method.String(), url.String(), body)
	if err != nil {
		return nil, err
	}

	r.Header = parseHeaders(in)

	return r, nil
}

type header struct {
	gohttp.Header
	status int
}

func newHeader() header {
	return header{
		Header: make(gohttp.Header),
		status: 200,
	}
}

type response[
	Headers http.Fields[F, N, V, HE], Trailers http.Fields[F, N, V, HE],
	Response any,
	F any, N http.FieldName, V http.FieldValue, HE http.HeaderError,
	OutgoingResponse http.OutgoingResponse[
		Headers, OutgoingBody,
		F, N, V, HE,
		O, I, SE, P,
	],
	OutgoingBody http.OutgoingBody[O, I, SE, P],
	O io.OutputStream[I, SE, P], I io.InputStream[SE, P], SE io.StreamError, P io.Pollable,
	ErrorCodeShape http.ErrorCodeShape, ErrorCode http.ErrorCode,
] struct {
	status              int
	header              header
	body                bytes.Buffer
	newHeaders          func() Headers
	newOutgoingResponse func(Headers) OutgoingResponse
	responseOutparamSet func(Response, cm.Result[ErrorCodeShape, OutgoingResponse, ErrorCode])
	outgoingBodyFinish  func(OutgoingBody, cm.Option[Trailers]) cm.Result[ErrorCode, struct{}, ErrorCode]
}

func newResponse[
	Headers http.Fields[F, N, V, HE], Trailers http.Fields[F, N, V, HE],
	Response any,
	F any, N http.FieldName, V http.FieldValue, HE http.HeaderError,
	OutgoingResponse http.OutgoingResponse[
		Headers, OutgoingBody,
		F, N, V, HE,
		O, I, SE, P,
	],
	OutgoingBody http.OutgoingBody[O, I, SE, P],
	O io.OutputStream[I, SE, P], I io.InputStream[SE, P], SE io.StreamError, P io.Pollable,
	ErrorCodeShape http.ErrorCodeShape, ErrorCode http.ErrorCode,
](
	newHeaders func() Headers,
	newOutgoingResponse func(Headers) OutgoingResponse,
	responseOutparamSet func(Response, cm.Result[ErrorCodeShape, OutgoingResponse, ErrorCode]),
	outgoingBodyFinish func(OutgoingBody, cm.Option[Trailers]) cm.Result[ErrorCode, struct{}, ErrorCode],
) *response[
	Headers, Trailers, Response,
	F, N, V, HE,
	OutgoingResponse, OutgoingBody,
	O, I, SE, P,
	ErrorCodeShape, ErrorCode,
] {
	return &response[
		Headers, Trailers, Response,
		F, N, V, HE,
		OutgoingResponse, OutgoingBody,
		O, I, SE, P,
		ErrorCodeShape, ErrorCode,
	]{
		header:              newHeader(),
		newHeaders:          newHeaders,
		newOutgoingResponse: newOutgoingResponse,
		responseOutparamSet: responseOutparamSet,
		outgoingBodyFinish:  outgoingBodyFinish,
	}
}

func (r *response[_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]) Header() gohttp.Header {
	return r.header.Header
}

func (r *response[_,_,  _, _, _, _, _, _, _, _, _, _, _, _, _]) Write(b []byte) (int, error) {
	r.body.Write(b)
	return len(b), nil
}

func (r *response[_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]) WriteHeader(statusCode int) {
	r.header.status = statusCode
}

func (r *response[_, Trailers, Response,
	_, FieldName, FieldValue, _,
	OutgoingResponse, _, _, _, _, _,
	ErrorCodeShape, ErrorCode,
]) flush(out Response) {
	headers := r.newHeaders()
	for k, vs := range r.header.Header {
		if vs == nil {
			continue
		}
		for _, v := range vs {
			headers.Append(FieldName(k), FieldValue(cm.ToList([]uint8(v))))
		}
	}

	w := r.newOutgoingResponse(headers)
	w.SetStatusCode(http.StatusCode(r.header.status))

	defer r.responseOutparamSet(out, cm.OK[cm.Result[ErrorCodeShape, OutgoingResponse, ErrorCode]](w))

	body, err := wit.HandleResult(w.Body())
	if err != nil {
		panic(fmt.Errorf("failed to get outgoing body: %s", err))
	}
	defer r.outgoingBodyFinish(*body, cm.None[Trailers]())

	output, err := wit.HandleResult((*body).Write())
	if err != nil {
		panic(fmt.Errorf("failed to write body: %s", err))
	}
	defer (*output).ResourceDrop()

	(*output).Write(cm.ToList(r.body.Bytes()))
}
