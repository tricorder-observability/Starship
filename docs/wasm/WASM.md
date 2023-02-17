# WASM Pre-study

## Original scope: web application

WebAssembly (WASM) is a portable binary instruction format which
was originally created for the web, it has a lot of use cases think about
image/video editing, games, VR, simulations and more.

* meant to execute heavy and intensive compute task of a web applications.
* can be compiled from other programming languages like C/C++, C#, Rust, Go and many more.

## New scope: system programming

WASM is starting to show up outside the browser,

* Relying on [WebAssembly System Interface (WASI)](https://github.com/WebAssembly/WASI);
* WASI is a modular collection of standardized APIs;

## Golang + WASM

WebAssembly can be compiled by Golang itself or with TinyGo.

* As with Go 1.19, it still does not support WASI, you need to use TinyGo if you wish to compile to WASI;

    TinyGo uses LLVM internally instead of emitting C, which hopefully leads to
    smaller and more efficient code and certainly leads to more flexibility.

    Goal:

    * Good CGo support, with no more overhead than a regular function call.
    * Support most standard library packages and compile most Go code without modification.

    Non-goals:

    * Be able to compile every Go program out there.

* Or, try [wasmtime-go](../demo/wasmtime-go/).

### Compile to JS

Examples:

1. [Golang WebAssembly](https://binx.io/2022/04/22/golang-webassembly/), 2022
1. [WebAssembly with Golang by scratch](https://itnext.io/webassemply-with-golang-by-scratch-e05ec5230558), 2021

### Compile to WASI with tinygo

Examples:

1. [Golang to WASI #PART1](https://www.wasm.builders/jennifer/golang-to-wasi-part1-il)
1. [Go in WebAssembly](https://www.fermyon.com/wasm-languages/go-lang), with golang 1.17
1. [WASI Hello World](https://wasmbyexample.dev/examples/wasi-hello-world/wasi-hello-world.go.en-us.html), 2021

### Compile with `wasmtime-go`

Package [wasmtime](https://pkg.go.dev/github.com/bytecodealliance/wasmtime-go)
is a WebAssembly runtime for Go powered by Wasmtime.

This package provides everything necessary to compile and execute WebAssembly
modules as part of a Go program. Wasmtime is a JIT compiler written in Rust,
and can be found at https://github.com/bytecodealliance/wasmtime. This package
is a binding to the C API provided by Wasmtime.
