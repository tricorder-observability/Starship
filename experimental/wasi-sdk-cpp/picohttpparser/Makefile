all: wasm wasm-lib wasm-lib2 wasm-lib3 wat go-wrapper

native:
	@rm -f ./pico
	clang -Wall -Wextra -Werror picohttpparser.c main.c -o pico

wasm:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -Wall -Wextra -Werror \
	  -Wl,--export=get_result_buf \
	  -Wl,--export=get_result_count \
	  -Wl,--export=get_result_struct_size \
	  picohttpparser.c main.c -o pico.wasm

wasm-lib:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -Wl,--no-entry \
	  -nostartfiles \
	  -Wall -Wextra -Werror \
	  -Wl,--allow-undefined \
	  -Wl,--export=allocate_input_output_bufs \
	  -Wl,--export=free_input_output_bufs \
	  -Wl,--export=pico_parse_events \
	  -Wl,--export=get_input_buf \
	  -Wl,--export=get_output_buf \
	  -Wl,--export=get_output_item_count \
	  -Wl,--export=get_output_item_size \
	  picohttpparser.c libpico.c -o libpico.wasm

wasm-lib2:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -nostartfiles \
	  -Wall -Wextra -Werror \
	  -Wl,--no-entry \
	  -Wl,--allow-undefined \
	  -Wl,--export=allocate_input_output_bufs \
	  -Wl,--export=free_input_output_bufs \
	  -Wl,--export=pico_parse_events \
	  -Wl,--export=get_output_buf \
	  -Wl,--export=get_output_item_count \
	  -Wl,--export=get_output_item_size \
	  picohttpparser.c libpico2.c -o libpico2.wasm

wasm-lib3:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -Wl,--no-entry \
	  -nostartfiles \
	  -Wall -Wextra -Werror \
	  -Wl,--allow-undefined \
	  -Wl,--export=pico_parse_events \
	  -Wl,--export=get_output_item_count \
	  -Wl,--export=get_output_item_size \
	  picohttpparser.c libpico3.c -o libpico3.wasm
	  # --target=wasm64 # wasi-sdk-17.0: only supports wasm32, doesn't yet support wasm64

wat:
	/opt/wabt-1.0.31/bin/wasm2wat pico.wasm > pico.wat

go-wrapper:
	go build main-wasmedge.go
	go build main-wasmtime.go
	go build main-wasmtime-libpico.go
	go build main-wasmtime-libpico2.go
	go build main-wasmtime-libpico3.go

clean:
	@rm -f *.wasm *.wat pico main-wasmedge main-wasmtime main-wasmtime-libpico main-wasmtime-libpico2
