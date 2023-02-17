all: wasm

native-bin:
	clang -o test test.c tlv_box.c key_list.c

wasm:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -Wall -Wextra -Werror \
	  tlv_box.c key_list.c test.c -o test.wasm

clean:
	@rm -f *.wasm *.wat test
