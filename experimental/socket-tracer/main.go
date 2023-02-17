// Copyright 2022 TriCorder
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/signal"

	bpf "github.com/iovisor/gobpf/bcc"
)

import "C"

// UniqPid stands for unique pid.
// See comments for struct u_pid_t {} in the BPF code for more explanations.
type UniqPid struct {
	pid            uint32 // PID or TGID: We use tgid in kernel-space, pid in user-space.
	startTimeTicks uint64
}

type connID struct {
	// The unique identifier of the pid/tgid.
	upid UniqPid
	// The file descriptor to the opened network connection.
	fd int32
	// Unique id of the conn_id (timestamp).
	tsid uint64
}

type enum uint32

type socketDataEventAttr struct {
	// The timestamp when syscall completed (return probe was triggered).
	timestamp_ns uint64

	// Connection identifier (PID, FD, etc.).
	conn_id connID

	// The protocol of traffic on the connection (HTTP, MySQL, etc.).
	protocol enum

	// The server-client role.
	role enum

	// The type of the actual data that the msg field encodes, which is used by the caller
	// to determine how to interpret the data.
	direction enum

	// Whether the traffic was collected from an encrypted channel.
	ssl bool

	// Represents the syscall or function that produces this event.
	source_fn enum

	// A 0-based position number for this event on the connection, in terms of byte position.
	// The position is for the first byte of this message.
	// Note that write/send have separate sequences than read/recv.
	pos uint64

	// The size of the original message. We use this to truncate msg field to minimize the amount
	// of data being transferred.
	msg_size uint32

	// The amount of data actually being sent to user space. This may be less than msg_size if
	// data had to be truncated, or if the data was stripped because we only want to send metadata
	// (e.g. if the connection data tracking has been disabled).
	msg_buf_size uint32

	// Whether to prepend length header to the buffer for messages first inferred as Kafka. MySQL
	// may also use this in this future.
	// See infer_kafka_message in protocol_inference.h for details.
	prepend_length_header bool
	length_header         uint32
}

type socketDataEvent struct {
	Pid         uint32
	Uid         uint32
	Gid         uint32
	ReturnValue int32
	Filename    [30702]byte // 30KB, defined along with the BPF code
}

type socketControlEvent struct {
	control_event_type_t enum
	timestamp_ns         uint64
	conn_id_t            connID

	// Represents the syscall or function that produces this event.
	source_fn enum

	// union {
	//   struct conn_event_t open;
	//   struct close_event_t close;
	// };
	close_event_t_wr_bytes int64
	close_event_t_rd_bytes int64
}

func main() {
	raw, err := os.ReadFile("./socket_trace.c")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open BPF source file: %s\n", err)
		os.Exit(1)
	}

	source := string(raw)
	// fmt.Printf("%s\n", source)

	m := bpf.NewModule(source, []string{})
	defer m.Close()

	// Attach BPF progs
	syscallName := bpf.GetSyscallFnName("connect")
	attachProbe(m, syscallName, "syscall__probe_entry_connect")
	attachKretprobe(m, syscallName, "syscall__probe_ret_connect")

	syscallName = bpf.GetSyscallFnName("close")
	attachProbe(m, syscallName, "syscall__probe_entry_close")
	attachKretprobe(m, syscallName, "syscall__probe_ret_close")

	// Open BPF map
	table := bpf.NewTable(m.TableId("socket_control_events"), m)
	channel := make(chan []byte)
	perfMap, err := bpf.InitPerfMap(table, channel, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		fmt.Printf("Reading perf buffers\n")

		// var event socketControlEvent
		for {
			data := <-channel
			// err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &event)
			// if err != nil {
			// 	fmt.Printf("failed to decode received data: %s\n", err)
			// 	continue
			// }
			fmt.Printf("Received event: %+v\n", string(data))
			// filename := (*C.char)(unsafe.Pointer(&event.Filename))
			// fmt.Printf("uid %d gid %d pid %d called fchownat(2) on %s (return value: %d)\n",
			// 	event.Uid, event.Gid, event.Pid, C.GoString(filename), event.ReturnValue)
		}
	}()

	perfMap.Start()
	<-sig
	perfMap.Stop()
}

func attachProbe(m *bpf.Module, syscallName, probeName string) {
	probe, err := m.LoadKprobe(probeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load %s: %s\n", probeName, err)
		os.Exit(1)
	}

	fmt.Printf("Probe loaded: id %+v, name %+v\n", probe, probeName)

	// passing -1 for maxActive signifies to use the default
	// according to the kernel kprobes documentation
	if err := m.AttachKprobe(syscallName, probe, -1); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach %s: %s\n", probeName, err)
		os.Exit(1)
	}
}

func attachKretprobe(m *bpf.Module, syscallName, probeName string) {
	probe, err := m.LoadKprobe(probeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load %s: %s\n", probeName, err)
		os.Exit(1)
	}

	fmt.Printf("Probe loaded: id %+v, name %+v\n", probe, probeName)

	// passing -1 for maxActive signifies to use the default
	// according to the kernel kprobes documentation
	if err := m.AttachKretprobe(syscallName, probe, -1); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach %s: %s\n", probeName, err)
		os.Exit(1)
	}
}
