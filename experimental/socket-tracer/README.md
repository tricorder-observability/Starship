socket-tracer demo
==================

This program attaches several BPF probes to socket related system calls
in kernel, collects and send the related events to a userspace agent.
The latter (written in golnag with gobpf) is responsible for reading,
compiling, loading and attaching the BPF program, as well as reading the
events and processing (now just printing) them.

# Prerquisites

* It is assumed that the local node has compiled and installed bcc from source
  with guide [Build BCC on Ubuntu 20.04](../doc/build-bcc).
* golang needs to be installed.

# Build

```shell
$ go version
go version go1.19.3 linux/amd64

$ make
go build -o socket-tracer main.go

$ file ./socket-tracer
./socket-tracer: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, BuildID[sha1]=3a3673bb6efe1c029a78b9729cbb889f740bacc1, for GNU/Linux 3.2.0, with debug_info, not stripped
```

# Test run

Star the userspace program:

```shell
$ sudo ./socket-tracer
Probe loaded: id 23, name syscall__probe_entry_connect
Probe loaded: id 25, name syscall__probe_ret_connect
Probe loaded: id 27, name syscall__probe_entry_close
Probe loaded: id 29, name syscall__probe_ret_close
Reading perf buffers
```

If no events show up, you can create some manually:

```shell
$ for n in {1..5}; do curl google.com; sleep 1; done
```

The output of the userspace agent:

```shell
$ sudo ./socket-tracer
Probe loaded: id 23, name syscall__probe_entry_connect
Probe loaded: id 25, name syscall__probe_ret_connect
Probe loaded: id 27, name syscall__probe_entry_close
Probe loaded: id 29, name syscall__probe_ret_close
Reading perf buffers
Received event: xxxxx (meaningless texts as I haven't decoded them, just print out with string(data)
Received event: xxxxx
Received event: ...
```

# Debug

BPF debug message is print to kernel tracing facility:

```shell
$ sudo tail /sys/kernel/debug/tracing/trace -n 20
   socket-tracer-46568   [014] d...1 82179.321511: bpf_trace_printk: Enter syscall__probe_ret_close, pid 199990857217512
   socket-tracer-46568   [014] d...1 82179.321512: bpf_trace_printk: Close args != NULL, 0
   socket-tracer-46568   [014] d...1 82179.321512: bpf_trace_printk: Enter process_syscall_close, 88
   socket-tracer-46568   [014] d...1 82179.321513: bpf_trace_printk: Exit before submit: conn_info === NULL, 0
   socket-tracer-46568   [013] d...1 82179.350703: bpf_trace_printk: Enter syscall__probe_entry_close, pid 199990857217512
   socket-tracer-46568   [013] d...1 82179.350716: bpf_trace_printk: Enter syscall__probe_entry_close, pid 199990857217512
```

Refer to [BPF 进阶笔记（四）：调试 BPF 程序](http://arthurchiao.art/blog/bpf-advanced-notes-4-zh) for more information about BPF debugging.

`bpftool` could also be used:

```shell
$ sudo bpftool prog show
...
336: kprobe  name syscall__probe_  tag cebbaaba0b6c9667  gpl
        loaded_at 2022-11-26T11:22:06+0000  uid 0
        xlated 528B  jited 299B  memlock 4096B  map_ids 424
        btf_id 66
337: kprobe  name syscall__probe_  tag 63057a207433f44e  gpl
        loaded_at 2022-11-26T11:22:06+0000  uid 0
        xlated 2192B  jited 1230B  memlock 4096B  map_ids 424,420,430,434,421,433
        btf_id 66

$ sudo bpftool map show
...
```

# Cleanup

Just stop the userspace program, then all resources will be released (as we didn't
pin anything to bpffs).

Remove any compilation file and the binary:

```shell
$ make clean
```

# TODO

1. Refactor BPF code, redefine some important event data structures (so they can be strictly matched with their userspace golang structs)
2. Decode events in userspace program
3. Attach more probes to kernel and collect the events
4. Test the L4/L7 processing BPF code
5. Cleanup
