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
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iovisor/gobpf/bcc"

	linux_headers "github.com/tricorder/src/agent/ebpf/bcc/linux-headers"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	testutils "github.com/tricorder/src/testing/bazel"
)

const code string = `
#include <linux/ptrace.h>
BPF_PERF_OUTPUT(events);
int sample_probe(struct bpf_perf_event_data* ctx) {
	const char word[] = "hello world";
	bpf_trace_printk("length=%d\n", sizeof(word));
	events.perf_submit(ctx, (void*)word, sizeof(word));
  return 0;
}
`

// Tests that AttachPerfEvent works as expected.
func TestAttachPerfEvent(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

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

	perfBuf, err := m.newPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

// Tests that attachKProbe and non syscall works as expected.
func TestAttachKprobe(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	assert.Nil(linux_headers.Init())

	m, err := newModule(code)
	require.Nil(err)
	defer m.Close()

	err = m.attachKProbe(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_KPROBE,
		Target: "ip_rcv",
		Entry:  "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err := m.newPerfBuffer("events")
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

	perfBuf, err = m.newPerfBuffer("events")
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

	assert.Nil(linux_headers.Init())

	m, err := newModule(code)
	require.Nil(err)
	defer m.Close()

	err = m.attachSyscallProbe(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_SYSCALL_PROBE,
		Target: "read",
		Entry:  "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err := m.newPerfBuffer("events")
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

	perfBuf, err = m.newPerfBuffer("events")
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

	assert.Nil(linux_headers.Init())

	m, err := newModule(code)
	require.Nil(err)
	defer m.Close()

	err = m.attachTracepoint(&ebpfpb.ProbeSpec{
		Type:   ebpfpb.ProbeSpec_TRACEPOINT,
		Target: "syscalls:sys_exit_read",
		Entry:  "sample_probe",
	})
	assert.Nil(err)

	perfBuf, err := m.newPerfBuffer("events")
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

	assert.Nil(linux_headers.Init())

	m, err := newModule(code)
	require.Nil(err)
	defer m.Close()

	const sampleUPROBEPath = "src/agent/ebpf/bcc/programs/test_uprobe"
	testBin := testutils.TestFilePath(sampleUPROBEPath)

	err = m.attachUProbe(&ebpfpb.ProbeSpec{
		Type:       ebpfpb.ProbeSpec_UPROBE,
		Target:     "main.sum",
		Entry:      "sample_probe",
		BinaryPath: testBin,
	})
	assert.Nil(err)

	perfBuf, err := m.newPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	cmd := exec.Command(testBin)
	assert.Nil(cmd.Run())

	time.Sleep(1 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()

	// return probe
	// golang uretprobe maybe happen some error
	// https://github.com/golang/go/issues/22008
	err = m.attachUProbe(&ebpfpb.ProbeSpec{
		Type:       ebpfpb.ProbeSpec_UPROBE,
		Target:     "main.sum",
		Return:     "sample_probe",
		BinaryPath: testBin,
	})
	assert.Nil(err)

	perfBuf, err = m.newPerfBuffer("events")
	require.Nil(err)
	perfBuf.Start()

	cmd = exec.Command(testBin)
	assert.Nil(cmd.Run())

	time.Sleep(1 * time.Second)
	bytesSlice = perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}

// Tests that the vanilla gobpf's BCC Golang binding APIs produce no extra null chars.
func TestDemoVanillaGoBPFAPI(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	const sampleJSONBPFCPath = "modules/sample_json/sample_json.bcc.c"
	bccCode, err := testutils.ReadTestFile(sampleJSONBPFCPath)
	require.Nil(err)

	m := bcc.NewModule(bccCode, []string{})
	defer m.Close()

	probeFD, err := m.LoadPerfEvent("sample_json")
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
		assert.Equal(`{"name":"John", "age":30}`+"\x00\x00\x00", string(data))
	}
	perfMap.Stop()
}
