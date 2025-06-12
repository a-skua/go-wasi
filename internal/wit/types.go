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
