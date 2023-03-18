all: copy_input_to_output.wasm

WASI_SDK_PATH := /opt/tricorder/wasm/wasi-sdk
WASI_CLANG := $(WASI_SDK_PATH)/bin/clang
CFLAGS := --sysroot=$(WASI_SDK_PATH)/share/wasi-sysroot -Wall -Wextra -Wl,--export-all

copy_input_to_output.wasm:
	$(WASI_CLANG) $(CFLAGS) copy_input_to_output.c -o $@

clean:
	rm -f *.wasm
