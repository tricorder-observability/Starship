# Common

Common libraries for writing WASM code compatible with Starship.  The north star
is to minimize such libraries. But at the beginning, we need a lot of
hand-written code in WASM to help them running on Tricorder's eBPF+WASM
infrastructure.

* `io.h`: Provides the APIs to allocate the input and output memory buffer
  inside the WASM runtime, and let the userspace to get the pointers.
* `cjson.{h,cc}`: Provides APIs to process JSON-formatted data in WASM C guest
  code.
