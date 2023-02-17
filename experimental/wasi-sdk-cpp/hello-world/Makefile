all:
	WASI_SDK_PATH=/opt/wasi-sdk-17.0 \
	/opt/wasi-sdk-17.0/bin/clang++ \
		hello.cc \
		-o hello.wasm

clean:
	@rm -f ./hello.wasm
