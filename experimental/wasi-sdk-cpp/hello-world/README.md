# Hello world

A simple example to illustrate how to compile a C++ program into WASI wasm and
execute it with command line tools.

## Build

```shell
~/starship/experimental/wasi-sdk-cpp/hello-world $ make
WASI_SDK_PATH=/opt/wasi-sdk-17.0 \
/opt/wasi-sdk-17.0/bin/clang++ \
    hello.cc \
    -o hello.wasm
```

The wasi-sdk provides a clang which is configured to target WASI and use the
WASI sysroot, so we can compile our program with `-o hello.wasm`. Now check
the generated target file:

```shell
~/starship/experimental/wasi-sdk-cpp/hello-world $ file ./hello.wam
hello.wasm: WebAssembly (wasm) binary module version 0x1 (MVP)
```

## Install `wasmtime` CLI

Tool `wasmtime` is needed to execute the generated `.wasm` binary file.

```shell
$ curl https://wasmtime.dev/install.sh -sSf | bash
```

## Run with `wasmtime` CLI

```shell
~/starship/experimental/wasi-sdk-cpp/hello-world $ wasmtime hello.wasm
Hello, World!
```
