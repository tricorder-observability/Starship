# Compiling pixie C++ parsing code to WASM

We've came up several ways to compile pixie C++ code into WASM and tried
some of them, but haven't succeeded yet.

## Fashion 1: resolve dependencies by hands

Steps to compile the http parsing module:

First, copy http parsing code to some place from pixie `src/stirling/source_connectors/socket_tracer/protocols/http/`;

Second, compile the code directly in `wasi-skd` way;

```shel
WASI_SDK_PATH=/opt/wasi-sdk-17.0 \
/opt/wasi-sdk-17.0/bin/clang++ \
	-I/usr/include \
	-I/usr/include/x86_64-linux-gnu/ \
	parse.cc \
    -o http-parser.wasm
```

which will use the clang/llvm bin/lib/headers provided by the sdk.

During this step, we encountered too many header dependencies problems,
and hundreds of pixie source files need to be modified to address
the relative path `include` problem, which is a big burden, so just gave up.

## Fashion 2: resolve dependencies by bazel in advance

In this way, we first try the vanilla step of pixie compilation:

```shell
$ cd pixie
$ bazel build src/stirling/source_connectors/socket_tracer/protocols/...:all
```

which can succeed if bazel works well. With step succeeded, all external dependencies
would have been resolved to the bazel cache in your home directory, such as
`/home/<user>/.cache/bazel/_bazel_<user>/f79d7e0fac5694b4cfb19cc3f07421dd/`.

With internal dependencies (`{.h, .cc}`) already in pixie repo itself, and
all external dependencies in the bazel cache, we could theoretically compile
the module like this:

```shell
# Yes, we're going to comiple our wasm target in vanilla pixie repo
$ cd ~/pixie/src/stirling/source_connectors/socket_tracer/protocols/http

$ Compile to wasm
$ WASI_SDK_PATH=/opt/wasi-sdk-17.0 \
  /opt/wasi-sdk-17.0/bin/clang++ \
	-I /home/<user>/pixie/ \
	-I/home/<user>/.cache/bazel/_bazel_<user>/f79d7e0fac5694b4cfb19cc3f07421dd/external/com_google_absl/ \
	-I/home/<user>/.cache/bazel/_bazel_<user>/f79d7e0fac5694b4cfb19cc3f07421dd/external/com_github_rlyeh_sole/ \
	-I/home/<user>/.cache/bazel/_bazel_<user>/f79d7e0fac5694b4cfb19cc3f07421dd/external/<other external deps>/ \
	-I/usr/include \
	-I/usr/include/x86_64-linux-gnu/ \
	parse.cc \
    -o http-parser.wasm
```

It worked indeed, until we encountered struct/function redefinition problems:

```shell
In file included from parse.cc:19:
In file included from /home/user/pixie/src/stirling/source_connectors/socket_tracer/protocols/http/parse.h:24:
In file included from /home/user/pixie/src/stirling/source_connectors/socket_tracer/protocols/common/interface.h:25:
In file included from /home/user/pixie/src/stirling/source_connectors/socket_tracer/bcc_bpf_intf/common.h:21:
In file included from /home/user/pixie/src/stirling/upid/upid.h:25:
In file included from /home/user/pixie/src/shared/upid/upid.h:27:
In file included from /home/user/.cache/bazel/_bazel_user/f79d7e0fac5694b4cfb19cc3f07421dd/external/com_github_rlyeh_sole/sole.hpp:178:
In file included from /usr/include/arpa/inet.h:22:
In file included from /usr/include/netinet/in.h:23:
In file included from /usr/include/sys/socket.h:26:
/usr/include/bits/types/struct_iovec.h:26:8: error: redefinition of 'iovec'
struct iovec
       ^
/opt/wasi-sdk-17.0/bin/../share/wasi-sysroot/include/__struct_iovec.h:7:8: note: previous definition is here
struct iovec {
```

This redefinition comes from the two direct inclusion of os-specific headers:

```shell
	-I/usr/include \
	-I/usr/include/x86_64-linux-gnu/ \
```

We haven't found a good way yet to solve this problem.

## Fashion 3: bazel + cross-compile

Haven't tried.

Pixie is working on cross-compile for arm arch.

https://publishing-project.rivendellweb.net/bazel-build-system-backend-code-even-more-choices/

## Fashion 4: bazel + wasm-sdk

Replace all the clang/llvm dependencies in pixie bazel configuration to the `wasi-sdk` equivalent.
Haven't knnow if it's possible.

## Fashion 5: vanilla clang/llvm

Haven't know if it's possible and how to do it.

Related:

1. https://github.com/ern0/howto-wasm-minimal
