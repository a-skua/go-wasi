package result

import (
	"fmt"

	"go.bytecodealliance.org/cm"
)

func HandleBool(r cm.BoolResult) bool {
	return r == cm.ResultOK
}

func Handle[Shape, OK, Err any](r cm.Result[Shape, OK, Err]) (OK, error) {
	ok, err, isErr := r.Result()
	if isErr {
		return ok, fmt.Errorf("error result: %v", err)
	}
	return ok, nil
}

func HandleErr[Shape, OK, Err any](r cm.Result[Shape, OK, Err], fn func(Err) error) (OK, error) {
	ok, err, isErr := r.Result()
	if isErr {
		return ok, fn(err)
	}
	return ok, nil
}

func Unwrap[Shape, OK, Err any](r cm.Result[Shape, OK, Err]) OK {
	ok, err, isErr := r.Result()
	if isErr {
		panic(fmt.Sprintf("result is an error: %v", err))
	}
	return ok
}

func UnwrapOr[Shape, OK, Err any](r cm.Result[Shape, OK, Err], defaultValue OK) OK {
	ok, _, isErr := r.Result()
	if isErr {
		return defaultValue
	}
	return ok
}
