package wit

import (
	"go.bytecodealliance.org/cm"
)

func HandleResult[Shape, OK, Err any](r cm.Result[Shape, OK, Err]) (*OK, *Err) {
	return r.OK(), r.Err()
}
