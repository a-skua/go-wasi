package option

import (
	"go.bytecodealliance.org/cm"
)

func Handle[T any](o cm.Option[T]) (zero T, _ bool) {
	if o.None() {
		return zero, false
	}
	return o.Value(), true
}

func Unwrap[T any](o cm.Option[T]) T {
	if o.None() {
		panic("option is None")
	}
	return o.Value()
}

func UnwrapOr[T any](o cm.Option[T], defaultValue T) T {
	if o.None() {
		return defaultValue
	}
	return o.Value()
}
