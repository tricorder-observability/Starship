# Run picohttpparser with wasm sandbox

Run `make wasm` to generate `pico.wasm`

`pico.wasm` is a WASI-style wasm binary file. It's built with WASI-enabled
compiler. It's source depends on system headers. But we actually can produce
a non-WASI binary file from picohttpparser.{h,cc}, by replacing stdio.h and
other headers. That can be compiled with the following command line:
```
clang --target=wasm32 --no-standard-libraries -Wl,--export-all -Wl,--no-entry -o pico.wasm picohttpparser.c
```

Compile picohttpparser (pure C) into wasm with wasi-sdk, and run the wasm
program with different WASM runtimes and their golang libraries,

* `wasmedge`, `wasmedge-go`
* `wasmtime`, `wasmtime-go`

Then retrieve the parsed results from sandbox memory.

## Design

In this section, we denote the golang application as "agent", and wasm sandbox adn "sandbox".

Input and output:

1. Agent -> Sandbox: via CLI arguments.

    Agent packs all the data that it wishes to send the sandbox as a C `argc, argv[]`
    argument list, then pass them to the sandbox.

    (According to my understanding) It is not possible to pass an agent side
    memory pointer to the sandbox, as this violates the security design of
    wasm.

2. Sandbox -> Agent: via sandbox memory.

    Wasm (C) program stores the parsing results to its own `malloc`ed memory
    (hooked by wasm runtime, if you know what I'm meaning), and the agent has
    the ability to retrieve this data. Several functions are exported by the
    wasm/c program to faciliate this data retrieval, see the code **and Makefile**.

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

## Debug

Some options for debugging and trouble shooting (save for later reference):

```shell
	# WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	# 			  -nostartfiles \
	# 			  -Wl,--no-entry \
	# 			  -Wl,--strip-all \
	# 			  -Wl,--export-all \
	# 			  -Wl,--import-memory \
	# 			  -fvisibility=hidden \
	#   -Wall -Wextra -Werror picohttpparser.c main.c -o pico.wasm

	# WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	# 			  -nostartfiles \
	# 			  -Wl,--no-entry \
	# 			  -Wl,--export=main2 \
	#   -Wall -Wextra -Werror picohttpparser.c main.c -o pico.wasm


	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
      -Wl,--trace-symbol=__wasi_fd_fdstat_get \
	  -Wall -Wextra -Werror picohttpparser.c main.c -o pico-main.wasm

	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
				  -Wl,--no-entry \
				  -nostartfiles \
	  -Wall -Wextra -Werror picohttpparser.c -o pico-lib.wasm
```

## Notes

* WASM function only accepts int32 and int64.
* Pass arguments as memory
* Compile code
* Return mutli value: https://github.com/bytecodealliance/wasmtime-go/blob/4b3a40c7a33dca1cfa4e06e6f40e06816f7d2421/example_test.go#L152
* AOT compiled binary, how to load and invoke.
  This helps improve performance
* https://github.com/bytecodealliance/wasmtime-go/blob/4b3a40c7a33dca1cfa4e06e6f40e06816f7d2421/example_test.go#L221
  can we use memory API.

