package bcc

import (
	"fmt"

	"github.com/iovisor/gobpf/bcc"
)

// Sleep 1 second waiting for data.
// PerfBuffer wraps a BCC perf buffer.
type PerfBuffer struct {
	Name       string
	bccTable   *bcc.Table
	channel    chan []byte
	bccPerfMap *bcc.PerfMap
}

func NewPerfBuffer(m *bcc.Module, name string) (*PerfBuffer, error) {
	res := new(PerfBuffer)
	res.Name = name
	res.bccTable = bcc.NewTable(m.TableId(name), m)
	// Create a buffered channel for bufferring data coming from perf buffer,
	// therefore we do not need a dedicated data buffer.
	res.channel = make(chan []byte, perfBufChanCap)

	var err error
	res.bccPerfMap, err = bcc.InitPerfMap(res.bccTable, res.channel, nil)

	if err != nil {
		return nil, fmt.Errorf("while creating PerfBuffer '%s', bcc.InitPerfMap() failed, error: %v", name, err)
	}
	return res, nil
}

func (perfBuf *PerfBuffer) Start() {
	perfBuf.bccPerfMap.Start()
}

func (perfBuf *PerfBuffer) Stop() {
	perfBuf.bccPerfMap.Stop()
}

// Poll returns all of the data currently in the perf buffer channel.
// Poll will not block if there is no data.
func (perfBuf *PerfBuffer) Poll() [][]byte {
	res := make([][]byte, 0)
	length := len(perfBuf.channel)
	for i := 0; i < length; i = i + 1 {
		item := <-perfBuf.channel
		res = append(res, item)
	}
	return res
}
