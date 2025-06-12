SRC := $(shell find . -name '*.go' -not -path './cmd/example/*' -not -path './internal/*')

.PROXY: example
example: cmd/example/http-proxy.wasm cmd/example/http-client.wasm

cmd/example/http-proxy.wasm: cmd/example/http-proxy/main.go world.wasm $(SRC)
	tinygo build -o $@ \
		--target=wasip2 --no-debug \
		--wit-package $(word 2, $^) \
		--wit-world http-proxy \
		$<

cmd/example/http-client.wasm: cmd/example/http-client/main.go world.wasm $(SRC)
	tinygo build -o $@ \
		--target=wasip2 --no-debug \
		--wit-package $(word 2, $^) \
		--wit-world http-client \
		$<

.PHONY: gen
gen: gen-http-proxy gen-http-client

.PHONY: gen-%
gen-%: world.wasm
	go tool wit-bindgen-go generate \
		--world $* \
		--out internal/gen $<

world.wasm: wit/world.wit
	wkg wit fetch
	wkg wit build -o $@

.PHONY: fmt
fmt:
	go fmt ./...
