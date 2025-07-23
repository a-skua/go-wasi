SRC := $(shell find ../.. -name '*.go' -not -path './cmd/example/*' -not -path './internal/gen/*')

.PROXY: example
example: http-proxy.wasm http-client.wasm

%.wasm: %/main.go world.wasm $(SRC)
	tinygo build -o $@ \
		--target=wasip2 --no-debug \
		--wit-package $(word 2, $^) \
		--wit-world $* \
		$<

world.wasm: wit/world.wit
	wkg wit fetch
	wkg wit build -o $@
