all: wasm wat go-wrapper

native:
	@rm -f ./pico
	clang -Wall -Wextra -Werror picohttpparser.c main.c -o pico

wasm:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -Wall -Wextra -Werror \
	  -Wl,--export=get_result_buf \
	  -Wl,--export=get_result_count \
	  -Wl,--export=get_result_struct_size \
	  -Wl,--export=phr_parse_request \
	  picohttpparser.c main.c -o pico.wasm

wat:
	/opt/wabt-1.0.31/bin/wasm2wat pico.wasm > pico.wat

go-wrapper:
	go build main-wasmedge.go
	go build main-wasmtime.go

clean:
	@rm -f *.wasm *.wat pico main-wasmedge main-wasmtime
