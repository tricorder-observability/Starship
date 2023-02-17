gobpf example
==============

This program attaches a simple BPF program to hook the `bash` `readline`
operations in kernel, and sends them as events to userspace. The userspace
agent is responsible for compiling, loading, attaching the BPF proram, as well
as receiving those events and processing (simple print) them.

# Prerquisites

* It is assumed that the local node has compiled and installed bcc from source
  with guide [Build BCC on Ubuntu 20.04](../doc/build-bcc).
* golang needs to be installed.

# Build

```shell
$ go version
go version go1.19.3 linux/amd64

$ go build bcc.go && file ./bcc
./bcc: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, BuildID[sha1]=cb5003205266fd5ff006a5f7f68233943d7cb3bd, for GNU/Linux 3.2.0, with debug_info, not stripped
```

# Test run

Star the userspace program:

```shell
$ sudo ./bcc
```

In another window, make some tests:

```shell
$ bash -c 'echo "foo"'
$ bash -c 'echo "bar"'
```

The output of the userspace agent:

```shell
$ sudo ./bcc
       PID      COMMAND
     27962      bash -c 'echo "foo"'
     27962      bash -c 'echo "bar"'
```

# Notes

Avoid to use `go run bcc.go` to run the example directly, as there may be
permission issues like below if you're not root user, and `sudo go run bcc.go`
won't solve them,

```shell
could not open bpf map: chowncall, error: Operation not permitted
panic: runtime error: invalid memory address or nil pointer dereference
        panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x4be314]
```
