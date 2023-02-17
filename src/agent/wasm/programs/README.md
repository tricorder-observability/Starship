# Programs

Store WASM guest programs in C used for demo and testing.
To build all .wasm files:
```
make clean
make
```

## struct-bindgen

Install ecc and struct-bindgen:

```sh
cd src/agent/wasm/libs/c
git clone https://github.com/eunomia-bpf/c-struct-bindgen --recursive
cd c-struct-bindgen
make
cp build/bin/Release/struct-bindgen ../struct-bindgen

# Download the latest version of ecc from eunomia-bpf releaess
# You could change the latest to a specily version
wget https://github.com/eunomia-bpf/eunomia-bpf/releases/latest/download/ecc \
    && chmod +x ./ecc
sudo apt install libclang-13 # may not need
sudo apt install libllvm15 # may not need
```

Generate WASM struct's header with
[ecc](https://github.com/eunomia-bpf/eunomia-bpf)
and
[struct-bindgen](https://github.com/eunomia-bpf/c-struct-bindgen):

```shell
./ecc event.h --header-only
./struct-bindgen event.bpf.o -p > struct-bindgen.h
clang-format -i struct-bindgen.h
```

One advantage of the approach in this PR is that it can handle pointers, since
now wasm32 only have 32-bit pointer. (Also, this approach is relatively common
in FFI-capable languages and bpftool). But pointer cannot be directly used
because wasm cannot directly read host memory. However, it could be store or
caculate offset in userspace, for example, print stack backtrace.

The `ecc` will call the clang to generate an bpf object file, which include the
BTF info of all the C structs and C types in the header. This include:

- Use `libclang` to find all struct definitions in the `AST` of the header file.
- For each struct, use clang `AST` rewrite to modified source code to generate
  BTF info.
- Use `clang` to do actural compile to generate the bpf object file.

This is the same as what the `bcc` frontend does, but ecc is a `AOT` compiler
rather than a `JIT` compiler, and `ecc` accepts standard C syntax file, rather
than a special C syntax file(call `BCC` style `C`).

The `struct-bindgen` will use the BTF info to generate the C header file, for
correct accessing the host struct in wasm.
