# GO + WASM Demo

This demo uses https://github.com/bytecodealliance/wasmtime-go golang library
to compile WASM source then run it on the fly.

Examples:

* hello-world: print hello world;
* hello-multi-args: passing multiple parameters to wasm program when calling it;
* externref-newmemory-wat.go: test externref and NewMemory with WAT code.

## Build

```shell
$ make
```

## Run

```shell
$ ./hello-world 
Hello from Go

$ ./hello-multi-args 
> 4 2
 0 1 2 3 4 5 6 7 8 9
```

## Clean

```shell
$ make clean
```
