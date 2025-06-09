.PHONY: gen
gen: gen_http_proxy

wit/world.wasm: wit/world.wit
	wkg wit fetch
	wkg wit build -o $@

.PHONY: gen_http_proxy
gen_http_proxy: wit/world.wasm
	go tool wit-bindgen-go generate --world http-proxy --out internal $<
