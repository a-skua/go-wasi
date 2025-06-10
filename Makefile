SRC := $(shell find . -name '*.go' -not -path './cmd/example/*' -not -path './internal/*')

.PROXY: example
example: cmd/example/http-proxy/main.wasm

cmd/example/%.wasm: cmd/example/%.go wit/world.wasm $(SRC)
	tinygo build -o $@ \
		--target=wasip2 --no-debug \
		--wit-package $(word 2, $^) \
		--wit-world http-proxy \
		$<

.PHONY: gen
gen: gen_http_proxy

.PHONY: gen_http_proxy
gen_http_proxy: wit/world.wasm
	go tool wit-bindgen-go generate --world http-proxy --out internal/gen $<

wit/world.wasm: wit/world.wit
	wkg wit fetch
	wkg wit build -o $@
