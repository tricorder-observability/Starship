# tlv_box

An easy-to-use TLV API in C.

Forked from https://github.com/Jhuster/TLV/tree/master/c,
fixed warnings, rewrote test, formatted code, and verified
compilation with wasi-sdk.

## Build

```shell
$ cd experimental/wasi-sdk-cpp/tlv
$ make
```

## Test

```shell
$ wasmtime ./test.wasm
Test encoding
Creating tlv box1
TLV box1 serialization successful, 106 bytes
Creating tlv box2
TLV box2 serialization successful, 114 bytes
Test encoding successful

Test decoding
Parse tlv box2 successful, 114bytes
Parse tlv box1 successful, 106 bytes
tlv_box_get_char successful x
tlv_box_get_short successful 2
tlv_box_get_int successful 3
tlv_box_get_long successful 4
tlv_box_get_float successful 5.670000
tlv_box_get_double successful 8.910000
tlv_box_get_string successful hello world!
tlv_box_get_bytes successful:  1-2-3-4-5-6-
Test decoding successful

Cleanup successful
```

## API

See `tlv_box.h` and `test.c`.
