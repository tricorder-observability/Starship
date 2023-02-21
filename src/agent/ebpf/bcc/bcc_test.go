// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package bcc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/agent/ebpf/bcc/linux_headers"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"

	"github.com/iovisor/gobpf/bcc"
)

const code string = `
#include <linux/ptrace.h>
BPF_PERF_OUTPUT(events);
int sample_probe(struct bpf_perf_event_data* ctx) {
	const char word[] = "hello world";
	bpf_trace_printk("submitting data ... \n");
	events.perf_submit(ctx, (void*)word, sizeof(word));
  return 0;
}
`

// Tests that AttachPerfEvent works as expected.
func TestAttachPerfEvent(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())

	m, err := newModule(code)
	require.Nil(err)
	defer m.Close()

	err = m.attachSampleProbe(&ebpfpb.ProbeSpec{
		Type:  ebpfpb.ProbeSpec_SAMPLE_PROBE,
		Entry: "sample_probe",

		SamplePeriodNanos: 100 * 1000 * 1000,
	})
	assert.Nil(err)

	perfBuf, err := m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

const probe_code string = `
#include <linux/skbuff.h>
#include <uapi/linux/ip.h>

BPF_PERF_OUTPUT(events);

int sample_probe(struct pt_regs *ctx, void *skb){
	const char word[] = "hello world";
	bpf_trace_printk("submitting data sample_probe... \n");
	events.perf_submit(ctx, (void*)word, sizeof(word));
}`

// Tests that attachKProbe and non syscall works as expected.
func TestAttachKprobe(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())

	m, err := newModule(probe_code)
	require.Nil(err)
	defer m.Close()

	err = m.attachKProbe(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_KPROBE,
		Target: "ip_rcv",
		Entry:  "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err := m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()

	// return probe
	err = m.attachKProbe(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_KPROBE,
		Target: "ip_rcv",
		Return: "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err = m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice = perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

// Tests that attachSyscallProbe works as expected.
func TestAttachSyscallProbe(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())

	m, err := newModule(probe_code)
	require.Nil(err)
	defer m.Close()

	err = m.attachSyscallProbe(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_SYSCALL_PROBE,
		Target: "read",
		Entry:  "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err := m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()

	// return probe
	err = m.attachSyscallProbe(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_SYSCALL_PROBE,
		Target: "read",
		Return: "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err = m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice = perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

// Tests that attachTracepoint works as expected.
func TestAttachTPProbe(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())

	m, err := newModule(probe_code)
	require.Nil(err)
	defer m.Close()

	err = m.attachTracepoint(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_TRACEPOINT,
		Target: "syscalls:sys_exit_read",
		Entry:  "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err := m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

// Tests that attachUprobe works as expected.
func TestAttachUProbe(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())

	m, err := newModule(probe_code)
	require.Nil(err)
	defer m.Close()

	err = m.attachUProbe(&ebpfpb.ProbeSpec{
		Type:       ebpfpb.ProbeSpec_UPROBE,
		Target:     "readline",
		Entry:      "sample_probe",
		BinaryPath: "/bin/bash",
	})
	assert.Nil(err)

	perfBuf, err := m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()

	// return probe
	err = m.attachUProbe(&ebpfpb.ProbeSpec{
		Type:       ebpfpb.ProbeSpec_UPROBE,
		Target:     "readline",
		Return:     "sample_probe",
		BinaryPath: "/bin/bash",
	})
	assert.Nil(err)

	perfBuf, err = m.NewPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice = perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

// Tests that the vanilla gobpf's BCC Golang binding APIs produce no extra null chars
func TestDemoVanillaGoBPFAPI(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := bcc.NewModule(code, []string{})
	defer m.Close()

	probeFD, err := m.LoadPerfEvent("sample_probe")
	require.Nil(err)

	err = m.AttachPerfEvent(1 /*evType*/, 0 /*evConfig*/, int(100000000), /*samplePeriod nanos*/
		ignoreSampleFreq, ignorePID, ignoreCPU, ignoreGroupFD, probeFD)
	require.Nil(err)

	table := bcc.NewTable(m.TableId("events"), m)

	channel := make(chan []byte, 1000)

	perfMap, err := bcc.InitPerfMap(table, channel, nil)
	require.Nil(err)

	perfMap.Start()
	for i := 0; i < 10; i++ {
		data := <-channel
		assert.Equal("hello world\x00", string(data))
	}
	perfMap.Stop()
}
