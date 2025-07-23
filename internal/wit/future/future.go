package future

import (
	"github.com/a-skua/go-wasi/internal/gen/wasi/io/poll"
)

type Future interface {
	Subscribe() poll.Pollable
}

func Wait[F Future](f F) {
	poll := f.Subscribe()
	defer poll.ResourceDrop()

	poll.Block()
}
