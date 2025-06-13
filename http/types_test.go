package http

import (
	"github.com/a-skua/go-wasi/internal/gen/wasi/http/types"
)

var (
	_ Scheme                                                       = &types.Scheme{}
	_ Method                                                       = &types.Method{}
	_ Fields[types.FieldName, types.FieldValue, types.HeaderError] = types.NewFields()
)
