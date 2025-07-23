.PHONY: examples
examples:
	@$(MAKE) -C examples

.PHONY: gen
gen: world.wasm
	@rm -rf internal/gen
	go tool wit-bindgen-go generate \
		--world wrapper \
		--out internal/gen $<

world.wasm: wit/world.wit
	wkg wit fetch
	wkg wit build -o $@

.PHONY: fmt
fmt:
	go fmt ./...
