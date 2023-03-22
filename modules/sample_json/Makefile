all: sample_json.wasm

WASI_SDK_PATH := /opt/tricorder/wasm/wasi-sdk
WASI_CLANG := $(WASI_SDK_PATH)/bin/clang
CFLAGS := --sysroot=$(WASI_SDK_PATH)/share/wasi-sysroot -Wall -Wextra -Wl,--export-all

sample_json.wasm:
	$(WASI_CLANG) $(CFLAGS) sample_json.c -o $@

clean:
	rm -f *.wasm
