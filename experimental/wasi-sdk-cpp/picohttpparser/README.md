# Run picohttpparser with wasm sandbox

Compile picohttpparser (pure C) into a wasm module with wasi-sdk, and run the wasm
program with different WASM runtimes and their golang libraries,

* `wasmedge`, `wasmedge-go`
* `wasmtime`, `wasmtime-go`

Then retrieve the parsed results from sandbox memory.

## Examples

Five examples provided, from naive to feature-rich:

1. main-wasmedge

    * compile pico into wasm executable
    * call wasm function with wasmedge runtime
    * pass input data to wasm function via argc/argv

1. main-wasmtime

    * compile pico into wasm executable
    * call wasm function with wasmtime runtime
    * pass input data to wasm function via argc/argv
    * get output data via wasm sandbox memory

1. main-wasmtime-libpico

    * compile pico into wasm library
    * call wasm function with wasmtime runtime
    * pass input data to wasm function via wasm sandbox memory
    * get output data via wasm sandbox memory

1. main-wasmtime-libpico2: **not working yet.**

    * Gogal: using `externref` in C code (that will be compiled to wasm).
    * Status: technically [impossible now](https://stackoverflow.com/questions/69457895/how-to-use-webassembly-reference-types-with-clang).

    Currently only possible if you're writing WebAssembly text code. See example_test.go in wasmtime-go repo.

    clang/llvm is developing a feature to support this:

    1. https://discourse.llvm.org/t/rfc-webassembly-tables-in-clang/62049
    2. https://reviews.llvm.org/D122215
    3. https://qiita.com/nokotan/items/5ca70b19818cd8776221

        This article is interesting, but I can't read Japanese, and google page translation failed :(

        One thing confirmed is that the `-mreference-types` option used in the
        article is not supported by our `/opt/clang-14/bin/clang` and `/opt/wasi-sdk-17.0/bin/clang`.
        Note that our wasi-sdk already has a higher verion (17.0) than the article's one (14.0).

        clang/llvm is also developing a feature to support WebAssembly table in native c/c++ code:
        https://reviews.llvm.org/D123510.

1. main-wasmtime-libpico3: **not working yet.**

    * Gogal: allocate memories for wasm module on host side, and transfer data with them;
    * Status: technically impossible now.

    We could successfully allocate memories for wasm module and pass memory addresses
    to it with host side function wrappers, but, on running the application, you'd
    got the following error:

    ```shell
    $ ./main-wasmtime-libpico3
    Go: creating new linker
    Go: creating new store
    Go: creating new memory for input data
    Go: instantiating sandbox
    Go: copying events into sandbox memory
    Go: calling wasm func pico_parse_events

    Go: input memory address: host side 0x7ff498000000
    WAMS: got input buffer 0x98000000
    Go: output memory address: host side 0x7ff298000000
    WAMS: got output buffer 0x98000000

    Received 3 http events, going to parse them
    2023/01/15 04:29:16 Call wasm func pico_parse_events failed: error while executing at wasm backtrace:
        0: 0x70c0 - <unknown>!strlen
        1: 0x358e - <unknown>!pico_parse_events
        2: 0x748e - <unknown>!pico_parse_events.command_export
    
    Caused by:
        wasm trap: out of bounds memory access
    ```

    Here were what happened in the behind:

    1. Host side successfully created an input memory and an output memory;
    2. WASM module successfully called the host side implemented getting memory address functions;
    3. On reading data from the input memory, the wasm sandbox trapped because of "out of bounds memory access".

    Reasons reasulted to this panic:

    1. Host side allocated wasm memories should not be passed to wasm module in this way.

        Host side is 64bit memory space, while wasm module is 32bit memory
        space (acceptting a i32 address pointer). Can we workaround this problem
        with wasmtime-go 64bit memory option? NO, both wasm & host side need to
        enable 64bit memory option, while wasi-sdk/clang doesn't
        support `wasm64` target yet.

    2. According to my current understanding, even if we could enable 64bit
       memory, and correctly passing a host allocated wasm memory to the wasm
       module, for our C program, it would still fail reading data, as these
       memories are non-default linear memories in the wasm module, while C
       programs still lack the ability to read non-default (or, additional) linear
       memories; currently, you can only do this with wasm text code (WAT).

## Build

```shell
$ cd experimental/wasi-sdk-cpp/picohttpparser
$ make
```

Three types of file will be generated:

1. WebAssembly binary format `.wasm`, for invoking from CLI or golang libraries;
2. WebAssembly text format `.wat`, for debugging;
3. Golang binaries: the simple applications.

## Test

### Invoke wasm program with golang wrapper

```shell
$ ./main-wasmtime $'GET /abc HTTP/1.1\r\nHost: tricorder.dev\r\nCookie: cookie\r\n\r\n'
$ ./main-wasmedge $'GET /abc HTTP/1.1\r\nHost: tricorder.dev\r\nCookie: cookie\r\n\r\n'
```

The outputs are much the same:

```shell
$ ./main-wasmtime $'GET /abcd HTTP/1.1\r\nHost: tricorder.dev\r\nCookie: cookie\r\n\r\n'
Go: creating new linker
Go: get and call function in wasm ...

picohttpparser started, going to process 1 events
Raw HTTP request: GET /abcd HTTP/1.1
Host: tricorder.dev
Cookie: cookie

Parse finished, ret=59, request information:
Method    : GET
Path      : /abcd
Header    : Host
          : tricorder.dev
Header    : Cookie
          : cookie
Copying return result to buffer, length 41
Copying return result finished

Go: returned from wasm sandbox main()
Go: getting get_result_buf
Go: getting get_result_len
Go: sandbox returned result: This is the return result from sandbox :)
```

Note that if http request is not provided to `main-wasmtime`, it will load some
built-in requests:

```
$ ./main-wasmtime
Go: creating new linker
Go: http request not provided through CLI, using built-in onesGo: get and call function in wasm ...

--------------------------- pico output (added by hand) ------------------------
Going to process 4 events
Parsed request information: # request 1
Method    : GET
Path      : /api/v1/bpf
Header    : Host
          : tricorder.dev
Header    : Cookie
          : cookie
Parsed request information: # request 2
...
--------------------------- pico output (added by hand) ------------------------

Go: calling into sandbox get_result_buf() ...
Go: calling into sandbox get_result_count() ...
Go: calling into sandbox get_result_struct_size() ...
Go: sandbox processed 3 events, result stored in buf addr 72352, result struct size 128

Go: parsed result decoded from sandbox output:
Method: GET, Path /api/v1/bpf
Method: PUT, Path /api/v2/wasm
Method: GET, Path /wp-content/uploads/2010/03/hello-kitty-darth-vader-pink.jpg
```

### Invoke wasm program with `wasmedge`/`wasmtime` CLI

Run the wasi file with `wasmtime` (or `wasmedge`) CLI:

```shell
$ wasmtime pico.wasm $'GET /abcd HTTP/1.1\r\nHost: tricorder.dev\r\nCookie: cookie\r\n\r\n'
picohttpparser started
Received raw HTTP request: GET /abcd HTTP/1.1 ...
Start parsing
Finish parsing, ret=59, parsed contents:
Method    : GET
Path      : /abcd
Header    : Host
          : tricorder.dev
Header    : Cookie
          : cookie
```

## Knowledge time: passing data to wasm module

In this section, denote the golang host application as "embedder", and wasm
sandbox program as "wasm module".

### Fashion 1: via CLI argc/argv

First, build the wasm module into an executable (instead of a library);

Then, embedder packs all the data that it would like to send to the asm program as a
standard C `argc, argv[]` arguments, and calls the program's `main` function
with those arguments. See main-wasmtime.go/main-wasmedge.go for examples.

### Fashion 2: via sandbox memory

First, build the wasm module into a library (instead of an executable);

Then, allocate a memory buffer in wasm sandbox's memory space, and export
the memory address to host side (embedder), serving as shared memory space.
Both the embedder and wasm module can read/write the buffer. But
the two sides should make a convention about the layout (data structure)
of the objects stored in the memory in order to correctly parse/unpack the data.

Usually, several functions should exported by the
wasm program to faciliate this data retrieval, see the `main-wasmtime.go`
or `main-wasmtime-libpico.go` and the **Makefile**.

### Fashion 3: via `externref`

Two frequently used ways to pass data from embedder to
wasm module are `NewMemory` and `externref`.

* NewMemory creates a new linear memory for a wasm sandbox, which can be used by
  the wasm module; passing the memory to wasm module needs the `table` feature.
* externref is a proposal of WebAssembly 2.0, which claims to allow passing any
  embedder object into wasm module directly.

Examples of these two can be found in wasmtime-go/`example_test.go`.
But the problem is that the receiving side (wasm) can only be written in
WebAssembly text (WAT) code. For any other languages, such as pure C in our case,
no automatic way exists; as mentioned above, clang/llvm supporting is still on the way.

WasmEdge provides a custom (and partial) externref implementation for Rust,
called [wasmedge-bingen](https://github.com/second-state/wasmedge-bindgen),
so if you're writing your wasm program with Rust, and running it with wasmedge runtime,
then you can pass or return objects of [several other types](https://github.com/second-state/wasmedge-bindgen/blob/main/bindgen/rust/macro)
other than the vanilla i32/i64 objects.
