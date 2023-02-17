all:
	rustup target add wasm32-wasi && rustc hello.rs --target wasm32-wasi
clean:
	@rm ./hello.wasm ./hello.wat
