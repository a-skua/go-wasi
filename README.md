# go-wasi
WASI light wrapper

## Examples

### Build

```sh
make
```

### http/proxy

```sh
wasmtime serve -S cli cmd/example/http-proxy/main.wasm
```

```sh
curl -i 'http://localhost:8080'
```

### Tools

1. [tinygo](https://tinygo.org/)
2. [wasmtime](https://wasmtime.dev/)
