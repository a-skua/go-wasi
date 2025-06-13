package io

import (
	"go.bytecodealliance.org/cm"
)

type StreamError interface {
}

type InputStream[E StreamError, P Pollable] interface {
	ResourceDrop()
	BlockingRead(len uint64) cm.Result[cm.List[uint8], cm.List[uint8], E]
	BlockingSkip(len uint64) cm.Result[uint64, uint64, E]
	Read(len uint64) cm.Result[cm.List[uint8], cm.List[uint8], E]
	Skip(len uint64) cm.Result[uint64, uint64, E]
	Subscribe() P
}

type OutputStream[I InputStream[E, P], E StreamError, P Pollable] interface {
	ResourceDrop()
	BlockingFlush() cm.Result[E, struct{}, E]
	BlokingSplice(I, uint64) cm.Result[uint64, uint64, E]
	BlockingWriteAndFlush(cm.List[uint8]) cm.Result[E, struct{}, E]
	BlockingWriteZeroesAndFlush(len uint64) cm.Result[E, struct{}, E]
	CheckWrite() cm.Result[uint64, uint64, E]
	Flush() cm.Result[E, struct{}, E]
	Splice(I, uint64) cm.Result[uint64, uint64, E]
	Subscribe() P
	WriteAndFlush(cm.List[uint8]) cm.Result[E, struct{}, E]
	Write(cm.List[uint8]) cm.Result[E, struct{}, E]
	WriteZeroes(len uint64) cm.Result[E, struct{}, E]
}
