package bcc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/agent/ebpf/bcc/linux_headers"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
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

		SamplePeriodNanos: 100,
	})
	assert.Nil(err)
}
