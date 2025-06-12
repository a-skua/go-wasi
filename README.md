# go-wasi
WASI light wrapper

## Examples

### Build

```sh
make
```

### e.g. http-proxy

```sh
wasmtime serve -S cli cmd/example/http-proxy.wasm
```

```sh
curl -i 'http://localhost:8080'
```

### e.g. http-client

```sh
wasmtime run -S http cmd/example/http-client.wasm
```

### Tools

1. [tinygo](https://tinygo.org/)
2. [wasmtime](https://wasmtime.dev/)
