package io

type Pollable interface {
	ResourceDrop()
	Block()
	Ready() bool
}
