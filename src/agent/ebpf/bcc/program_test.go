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

	"github.com/tricorder/src/agent/ebpf/bcc/linux_headers"
	commonpb "github.com/tricorder/src/pb/module/common"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
)

const bccCode string = `
#include <linux/ptrace.h>
BPF_PERF_OUTPUT(events);
int syscall__probe_entry_read(struct pt_regs* ctx, int fd, char* buf, size_t count) {
	const char word[] = "hello world";
	bpf_trace_printk("submitting data ... \n");
	events.perf_submit(ctx, (void*)word, sizeof(word));
  return 0;
}
`

func loadProgram(t *testing.T, code *ebpfpb.Program) *Program {
	prog, err := NewProgram(code)
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = prog.Init()
	if err != nil {
		t.Errorf("%v", err)
	}

	return prog
}

func TestLoadAndPollData(t *testing.T) {
	assert := assert.New(t)

	// init kernel headers
	assert.Nil(linux_headers.Init())
	ebpfProgram := ebpfpb.Program{
		Fmt:            commonpb.Format_TEXT,
		Lang:           commonpb.Lang_C,
		Code:           bccCode,
		PerfBufferName: "events",
		Probes: []*ebpfpb.ProbeSpec{
			{
				Type:   ebpfpb.ProbeSpec_SYSCALL_PROBE,
				Target: "read",
				Entry:  "syscall__probe_entry_read",
			},
		},
	}
	prog := loadProgram(t, &ebpfProgram)
	assert.Nil(prog.Init())

	// Sleep 1 second waiting for data.
	time.Sleep(time.Second)

	perfBufData := prog.Poll()
	assert.NotNil(perfBufData)

	bytes, found := perfBufData["events"]
	assert.True(found)
	assert.Equal("hello world\x00", string(bytes[0]))
	prog.Stop()
}
