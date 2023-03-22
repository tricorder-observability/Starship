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
	"fmt"

	"github.com/tricorder/src/utils/log"

	ebpfpb "github.com/tricorder/src/pb/module/ebpf"

	"github.com/tricorder/src/utils/pb"
)

// perfBufChanCap gives the capacity of the perf buffer channel.
var perfBufChanCap = 1000

// Program abstract a piece of eBPF program. Provides APIs for managing the program's
// whole lifetime, and APIs for interacting with the attached eBPF program.
// For example, polling perf buffer and get the data from eBPF program collected from inside Kernel.
type Program struct {
	mod  *module
	spec *ebpfpb.Program

	PerfBufferNames []string
	perfBuffers     []*PerfBuffer
}

func NewProgram(p *ebpfpb.Program) (*Program, error) {
	res := new(Program)

	m, err := newModule(p.Code)
	if err != nil {
		return nil, fmt.Errorf("while creating Program, failed to create BCC Module, error: %v", err)
	}
	res.mod = m
	res.spec = p
	res.PerfBufferNames = []string{p.PerfBufferName}
	return res, nil
}

func (p *Program) Init() error {
	for _, probe := range p.spec.Probes {
		log.Infof("Attaching probe: %s", pb.FormatOneLine(probe))
		if err := p.mod.attachProbe(probe); err != nil {
			return fmt.Errorf("failed to attach probe '%s', error: %v", probe, err)
		}
	}
	perfBuffer, err := p.mod.newPerfBuffer(p.spec.PerfBufferName)
	if err != nil {
		return fmt.Errorf("while initializing eBPF program, failed to create PerfBuffer, error: %v", err)
	}
	perfBuffer.Start()
	p.perfBuffers = append(p.perfBuffers, perfBuffer)
	return nil
}

func (p *Program) Poll() map[string][][]byte {
	res := make(map[string][][]byte)
	for _, buf := range p.perfBuffers {
		res[buf.Name] = buf.Poll()
	}
	return res
}

func (p *Program) Stop() {
	// TODO(yaxiong): Unload eBPF program, delete WASM etc.
	for _, perfBuf := range p.perfBuffers {
		perfBuf.Stop()
	}
	p.mod.Close()
}
