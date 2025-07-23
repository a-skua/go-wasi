SRC := $(shell find .. -name '*.go' -not -path '../examples/*' -not -path '../internal/gen/*')

.PHONY: examples
examples: http-proxy.wasm http-client.wasm print.wasm

%.wasm: %/main.go world.wasm $(SRC)
	tinygo build -o $@ \
		--target=wasip2 --no-debug \
		--wit-package $(word 2, $^) \
		--wit-world $* \
		$<

world.wasm: wit/world.wit
	wkg wit fetch
	wkg wit build -o $@
