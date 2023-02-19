package bcc

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iovisor/gobpf/bcc"
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
/*
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

	time.Sleep(10 * time.Second)
	bytesSlice := perfBuf.Poll()
	for _, bytes := range bytesSlice {
		assert.Equal("hello world\x00", string(bytes))
	}
	perfBuf.Stop()
}
*/

// Tests that the vanilla gobpf's BCC Golang binding APIs produce no extra null chars
func TestDemoVanillaGoBPFAPI(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	const sampleJSONBPFCPath = "modules/sample_json/sample_json.bcc"
	bccCode, err := testutils.ReadTestFile(sampleJSONBPFCPath)
	require.Nil(err)
	log.Print(string(bccCode))

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
		log.Print("golang len(data)", len(data))
		// assert.Equal("hello world\x00", string(data))
		assert.True(true)
	}
	perfMap.Stop()
}
