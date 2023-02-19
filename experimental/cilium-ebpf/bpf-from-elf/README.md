# Parse and Run BPF program from a given BPF ELF file

This piece of code demonstrates how to parse a user-provided ELF file.
With the BPF maps and programs information parsed from the ELF, a simple agent
loads and attaches the BPF programs (and creates BPF maps) to kernel dynamically.
If the name of the BPF map for sending events from kernel to userspace is given,
the agent will also create a corresponding perf/ring buffer, then polling data
from the buffer.

Note: the code is just for POC, so lots of error handlings are missing.

## Supported tracing types

Tested:

* Tracepoint
* Kprobe
* fentry/fexit

More can be added.

## Test

First, generate some BPF object files (ELF) for testing. Here we use an existing
project `libbpf-bootstrap`:

```shell
# Clone code and init submodule
$ git clone https://github.com/libbpf/libbpf-bootstrap.git
$ cd libbpf-bootstrap
$ git submodule update --init --recursive       # check out libbpf

# Make examples
$ cd ~/libbpf-bootstrap/examples/c
$ make -j4

# Check the generated files
$ ls ~/libbpf-bootstrap/examples/c/.output/
bootstrap.bpf.o   bpftool        kprobe.bpf.o   libbpf.a              minimal_legacy.skel.h  sockfilter.bpf.o   tc.o          uprobe.skel.h
bootstrap.o       fentry.bpf.o   kprobe.o       minimal.bpf.o         minimal.o              sockfilter.o       tc.skel.h     usdt.bpf.o
bootstrap.skel.h  fentry.o       kprobe.skel.h  minimal_legacy.bpf.o  minimal.skel.h         sockfilter.skel.h  uprobe.bpf.o  usdt.o
bpf               fentry.skel.h  libbpf         minimal_legacy.o      pkgconfig              tc.bpf.o           uprobe.o      usdt.skel.h
```

We'll use the above generated `*.bpf.o` files (BPF ELFs) as input.

Now use our cilium/ebpf based agent to parse, load and attach BPF programs
into kernel, then poll data from the BPF maps specified in the ELF:

```shell
# Build our agent
$ make
```

```shell
# Usage: ./main <BPF ELF file> <event map name>
# "rb" is the BPF map (ring buffer) name defined in bootstrap.bpf.c
$ ./main ~/libbpf-bootstrap/examples/c/.output/bootstrap.bpf.o "rb"
2023/02/19 09:17:47 Input ELF file: ~/libbpf-bootstrap/examples/c/.output/bootstrap.bpf.o
2023/02/19 09:17:47 Load CollectionSpec from ELF successful

Dumping ELF collection spec
BPF map: name rb, spec RingBuf(keySize=0, valueSize=0, maxEntries=262144, flags=0)
BPF map: name .rodata, spec Array(keySize=4, valueSize=8, maxEntries=1, flags=128)
BPF map: name exec_start, spec Hash(keySize=4, valueSize=8, maxEntries=8192, flags=0)
BPF program: name handle_exec, spec &{Name:handle_exec Type:TracePoint AttachType:None AttachTo:sched/sched_process_exec AttachTarget:<nil> SectionName:tp/sched/sched_process_exec Instructions:handle_exec:
	   ; int handle_exec(struct trace_event_raw_sched_process_exec *ctx)
	  0: MovReg dst: r6 src: r1
	   ; pid = bpf_get_current_pid_tgid() >> 32;
	  1: Call FnGetCurrentPidTgid
	   ; pid = bpf_get_current_pid_tgid() >> 32;
	  2: RShImm dst: r0 imm: 32
	   ...
	 62: Exit
 Flags:0 License:Dual BSD/GPL KernelVersion:0 ByteOrder:LittleEndian}
BPF program: name handle_exit, spec &{Name:handle_exit Type:TracePoint AttachType:None AttachTo:sched/sched_process_exit AttachTarget:<nil> SectionName:tp/sched/sched_process_exit Instructions:handle_exit:
	   ; id = bpf_get_current_pid_tgid();
	  0: Call FnGetCurrentPidTgid
	   ; pid = id >> 32;
	  1: MovReg dst: r1 src: r0
	  2: RShImm dst: r1 imm: 32
	   ; pid = id >> 32;
	  ...
	 83: MovImm dst: r0 imm: 0
	 84: Exit
 Flags:0 License:Dual BSD/GPL KernelVersion:0 ByteOrder:LittleEndian}
Dumping ELF collection spec done

Load BPF collection (programs/maps) into kernel successful
Dumping ELF collection
BPF map: name exec_start, spec {name:exec_start fd:0xc00024bc90 typ:1 keySize:4 valueSize:8 maxEntries:8192 flags:0 pinnedPath: fullValueSize:8}
BPF map: name rb, spec {name:rb fd:0xc00024bca0 typ:27 keySize:0 valueSize:0 maxEntries:262144 flags:0 pinnedPath: fullValueSize:0}
BPF map: name .rodata, spec {name:.rodata fd:0xc00024bd70 typ:2 keySize:4 valueSize:8 maxEntries:1 flags:128 pinnedPath: fullValueSize:8}
BPF program: name handle_exit, spec {VerifierLog: fd:0xc000976a68 name:handle_exit pinnedPath: typ:5}
BPF program: name handle_exec, spec {VerifierLog: fd:0xc000976300 name:handle_exec pinnedPath: typ:5}
Dumping ELF collection done

2023/02/19 09:17:47 Attaching BPF program TracePoint(handle_exec)#10 to handle_exec
2023/02/19 09:17:47 Attaching BPF program TracePoint(handle_exec)#10 to handle_exec: group sched, name sched_process_exec
2023/02/19 09:17:47 Attaching BPF program TracePoint(handle_exit)#11 to handle_exit
2023/02/19 09:17:47 Attaching BPF program TracePoint(handle_exit)#11 to handle_exit: group sched, name sched_process_exit
2023/02/19 09:17:47 Create ringbuffer for polling data from BPF map rb successful
2023/02/19 09:17:50 Got an event from BPF map, event size 168
2023/02/19 09:17:51 Got an event from BPF map, event size 168
2023/02/19 09:17:51 Got an event from BPF map, event size 168
2023/02/19 09:17:51 ...
2023/02/19 09:17:51 Got an event from BPF map, event size 168
2023/02/19 09:17:51 Received stop signal, notify goroutines to exit
2023/02/19 09:17:51 Received stop signal, close BPF link for program handle_exec
2023/02/19 09:17:51 Received stop signal, close polling reader for BPF map rb
2023/02/19 09:17:51 Received stop signal, close BPF link for program handle_exit
2023/02/19 09:17:51 received signal, exiting..
2023/02/19 09:17:52 Agent exited
```

Other tests:

```shell
$ ./main ~/libbpf-bootstrap/examples/c/.output/kprobe.bpf.o
$ ./main ~/libbpf-bootstrap/examples/c/.output/fentry.bpf.o
```

These two examples do not have data output to userspace, but only prints logs
to kernel debugging facility, you can check them with:

```shell
$ sudo cat /sys/kernel/debug/tracing/trace_pipe
```

## Limitations

Event data structure information is missing from ELF, so we can only get events
from kernel, but can't decode them in usespace (for debugging).
