all: struct_test.wasm

WASI_CLANG := WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang
CFLAGS := --sysroot=/opt/wasi-sdk/share/wasi-sysroot -Wall -Wextra -Wl,--export-all

struct_test.wasm:
	$(WASI_CLANG) $(CFLAGS) struct_test.c -o $@

clean:
	rm -f *.wasm
