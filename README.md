# go-wasi
WASI light wrapper

## Examples

### Build

```sh
make
```

### e.g. http-proxy

```sh
wasmtime serve -S cli examples/http-proxy.wasm
```

```sh
curl -i 'http://localhost:8080'
```

### e.g. http-client

```sh
wasmtime run -S http examples/http-client.wasm
```


### e.g. print

```sh
wasmtime run --dir=./examples::/ examples/print.wasm /Makefile /foo
```

### Tools

1. [tinygo](https://tinygo.org/)
2. [wasmtime](https://wasmtime.dev/)
