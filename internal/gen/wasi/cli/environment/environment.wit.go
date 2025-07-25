// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package environment represents the imported interface "wasi:cli/environment@0.2.5".
package environment

import (
	"go.bytecodealliance.org/cm"
)

// GetEnvironment represents the imported function "get-environment".
//
//	get-environment: func() -> list<tuple<string, string>>
//
//go:nosplit
func GetEnvironment() (result cm.List[[2]string]) {
	wasmimport_GetEnvironment(&result)
	return
}

// GetArguments represents the imported function "get-arguments".
//
//	get-arguments: func() -> list<string>
//
//go:nosplit
func GetArguments() (result cm.List[string]) {
	wasmimport_GetArguments(&result)
	return
}

// InitialCWD represents the imported function "initial-cwd".
//
//	initial-cwd: func() -> option<string>
//
//go:nosplit
func InitialCWD() (result cm.Option[string]) {
	wasmimport_InitialCWD(&result)
	return
}
