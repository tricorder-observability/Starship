# Demo

Install rust building environment:

```shell
$ curl https://sh.rustup.rs -sSf | sh
```

Compile:

```shell
$ make
rustup target add wasm32-wasi && rustc hello.rs --target wasm32-wasi
info: component 'rust-std' for target 'wasm32-wasi' is up to date

$ file hello.wasm 
hello.wasm: WebAssembly (wasm) binary module version 0x1 (MVP)
```

Install `wasmtime` if you would like to run it (I haven't tried this):

```shell
$ wasmtime hello.wasm
Hello, world!
```

## Converting binary format to text format

Download https://github.com/WebAssembly/wabt/releases tools, for ubuntu:

```shell
$ curl https://github.com/WebAssembly/wabt/releases/download/1.0.31/wabt-1.0.31-ubuntu.tar.gz
$ tar xvf wabt-1.0.31-ubuntu.tar.gz
$ mv wabt-1.0.31 ~
```

Now perform converting:

```shell
$ ~/wabt-1.0.31/bin/wasm2wat ./hello.wasm > hello.wat

$ head hello.wat 
(module
  (type (;0;) (func))
  (type (;1;) (func (param i32)))
  (type (;2;) (func (param i32) (result i64)))
  (type (;3;) (func (param i32 i32)))
  (type (;4;) (func (param i32) (result i32)))
  (type (;5;) (func (param i32 i32) (result i32)))
  (type (;6;) (func (param i32 i32 i32)))
  (type (;7;) (func (param i32 i32 i32) (result i32)))
  (type (;8;) (func (param i32 i32 i32 i32) (result i32)))

$ wc -l hello.wat
25139 hello.wat
```

## Summary
Rust code can be directly compiled into the wasm binary format (.wasm), which can be excuted by wasmtime runtime/CLI.

wabt is a WebAssembly Binary Toolkit, which includes many useful tools such as wasm2wat to conver the binary format into text format (WebAssembly source file).
