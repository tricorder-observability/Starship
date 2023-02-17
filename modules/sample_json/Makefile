all: copy_input_to_output.wasm

WASI_CLANG := WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang
CFLAGS := --sysroot=/opt/wasi-sdk/share/wasi-sysroot -Wall -Wextra -Wl,--export-all

copy_input_to_output.wasm:
	$(WASI_CLANG) $(CFLAGS) copy_input_to_output.c -o $@

clean:
	rm -f *.wasm
