package wit

import (
	"fmt"
	"go.bytecodealliance.org/cm"
)

func HandleResult[Shape, OK, Err any](r cm.Result[Shape, OK, Err]) (*OK, *Err) {
	return r.OK(), r.Err()
}

func UnwrapResult[Shape, OK, Err any](r cm.Result[Shape, OK, Err]) *OK {
	if r.IsErr() {
		panic(fmt.Sprintf("result is an error: %v", r.Err()))
	}
	return r.OK()
}

func HandleOption[T any](o cm.Option[T]) (value T, ok bool) {
	if o.None() {
		var zero T
		return zero, false
	}
	return o.Value(), true
}

func UnwrapOption[T any](o cm.Option[T]) T {
	if o.None() {
		panic("option is None")
	}
	return o.Value()
}

func UnwrapOptionOr[T any](o cm.Option[T], defaultValue T) T {
	if o.None() {
		return defaultValue
	}
	return o.Value()
}
