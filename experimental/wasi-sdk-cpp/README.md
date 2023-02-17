# Compiling C/C++ code to WASM with  `wasi-sdk`

SDK: github.com/WebAssembly/wasi-sdk

> WASI-enabled WebAssembly C/C++ toolchain.

## What's WASI?

From [wasmtime documentation](https://github.com/bytecodealliance/wasmtime/blob/main/docs/WASI-overview.md#wasi-webassembly-system-interface)
(highly recommended to read):

> WebAssembly System Interface, or WASI, is a new family of API's being designed
> by the Wasmtime project to propose as a standard engine-independent non-Web
> system-oriented API for WebAssembly.

WASI software architecture:

![](https://github.com/bytecodealliance/wasmtime/raw/main/docs/wasi-software-architecture.png)

## Install `wasi-sdk`

Just download the pre-built binary and unzip file to a place, such as `/opt/`
(recommended by official doc).

## Limitations

Accroding to
https://github.com/WebAssembly/wasi-sdk#notable-limitations:

1. Do not support C++ exceptions yet. C++ code is supported only with -fno-exceptions for now.
2. Do not support threads yet.
3. Do not support dynamic library yet.
4. Do nto support networking yet.

## Recommended readings

1. [WASI Overview](https://github.com/bytecodealliance/wasmtime/blob/main/docs/WASI-overview.md#wasi-webassembly-system-interface), wasmtime documentation
2. [`wasi-sdk` C/Rust tutorial](https://github.com/bytecodealliance/wasmtime/blob/main/docs/WASI-tutorial.md), wasmtime documentation
3. [wasmtime online documentation (book)](https://docs.wasmtime.dev/)

# compile with emsdk

## reference

- https://emscripten.org/docs/compiling/Building-Projects.html
- https://github.com/emscripten-core/emsdk

## install

```sh
# emscripten install
git clone https://github.com/emscripten-core/emsdk.git
cd emsdk
./emsdk update-tags
./emsdk install 3.1.7
./emsdk activate 3.1.7
source ./emsdk_env.sh
cd
```

## example: compile abseil

```sh
# abseil (optional)
git clone https://github.com/abseil/abseil-cpp
cd abseil-cpp
emcmake cmake -DCMAKE_CXX_STANDARD=17 "."
emmake make
cd
```

## possible working commands

You should compile picohttpparser, abseil-cpp, sole, gflags, magic_enum, protobuf... etc to wasm, or remove them from include headers. all linux include headers should be remove.

```sh
EMCC_ONLY_FORCED_STDLIBS=1 em++ -O3 -s STANDALONE_WASM=1 \
    -o parse.wasm parse.cc \
    -I/home/yunwei/pixie/ \
    -I/home/yunwei/deps/picohttpparser \
    -I/home/yunwei/deps/protobuf/third_party/abseil-cpp \
    -I/home/yunwei/deps/sole \
    -I/home/yunwei/deps/gflags/include \
    -I/home/yunwei/deps/protobuf/src \
    -I/home/yunwei/deps/magic_enum/include \
    -I/home/yunwei/deps/glog/ \
    -s TOTAL_STACK=4096 -s TOTAL_MEMORY=65536 \
    -s ERROR_ON_UNDEFINED_SYMBOLS=0
```

because emsdk has makefile, cmake support, it should be easier to use it.

## an example hacking patch for compile http decoder

```
From 4863c2bb52ae2d022fdaa34988cc7bc20f6282cd Mon Sep 17 00:00:00 2001
From: yunwei37 <1067852565@qq.com>
Date: Sun, 18 Dec 2022 15:11:33 +0800
Subject: [PATCH] patch to compile decoder

---
 .../protocols/http/body_decoder.h             | 27 ++++++++++++++++---
 src/stirling/utils/parse_state.h              | 13 ++-------
 2 files changed, 26 insertions(+), 14 deletions(-)

diff --git a/src/stirling/source_connectors/socket_tracer/protocols/http/body_decoder.h b/src/stirling/source_connectors/socket_tracer/protocols/http/body_decoder.h
index 635531a..a138049 100644
--- a/src/stirling/source_connectors/socket_tracer/protocols/http/body_decoder.h
+++ b/src/stirling/source_connectors/socket_tracer/protocols/http/body_decoder.h
@@ -19,11 +19,32 @@
 #pragma once
 
 #include <string>
-
+#include <iostream>
 #include "src/stirling/utils/parse_state.h"
+#include <absl/strings/substitute.h>
+
+#define DECLARE_VARIABLE(type, shorttype, name, tn) \
+  namespace fL##shorttype {                         \
+    extern type FLAGS_##name;           \
+  }                                                 \
+  using fL##shorttype::FLAGS_##name
+#define DEFINE_VARIABLE(type, shorttype, name, value, meaning, tn) \
+  namespace fL##shorttype {                                        \
+    type FLAGS_##name(value);                          \
+    char FLAGS_no##name;                                           \
+  }                                                                \
+  using fL##shorttype::FLAGS_##name
+
+// bool specialization
+#define DECLARE_bool(name) \
+  DECLARE_VARIABLE(bool, B, name, bool)
+#define DEFINE_bool(name, value, meaning) \
+  DEFINE_VARIABLE(bool, B, name, value, meaning, bool)
+
+#define LOG(INFO) std::cout
+#define DFATAL
+#define ERROR
 
-// Choose either the pico or custom implementation of the chunked HTTP body decoder.
-DECLARE_bool(use_pico_chunked_decoder);
 
 namespace px {
 namespace stirling {
diff --git a/src/stirling/utils/parse_state.h b/src/stirling/utils/parse_state.h
index 807c945..1c5adc0 100644
--- a/src/stirling/utils/parse_state.h
+++ b/src/stirling/utils/parse_state.h
@@ -21,10 +21,9 @@
 #include <string>
 #include <vector>
 
-#include "src/common/base/base.h"
-
 namespace px {
 namespace stirling {
+struct Status;
 
 enum class ParseState {
   kUnknown,
@@ -52,15 +51,7 @@ enum class ParseState {
   kSuccess,
 };
 
-inline ParseState TranslateStatus(const Status& status) {
-  if (error::IsNotFound(status) || error::IsResourceUnavailable(status)) {
-    return ParseState::kNeedsMoreData;
-  }
-  if (!status.ok()) {
-    return ParseState::kInvalid;
-  }
-  return ParseState::kSuccess;
-}
+inline ParseState TranslateStatus(const Status& status);
 
 }  // namespace stirling
 }  // namespace px
-- 
2.37.2


```
