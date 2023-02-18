package driver

import (
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tricorder/src/agent/ebpf/bcc"
	"github.com/tricorder/src/agent/wasm"
	"github.com/tricorder/src/utils/pg"

	modulepb "github.com/tricorder/src/pb/module"
)

// Module holds data about an eBPF+WASM module waiting for being deployed.
type Module struct {
	modulePB *modulepb.Module

	// An abstract of a BCC program, which provides interfaces to manage the whole lifetime
	// of the BCC program and performs various operation during it.
	ebpf *bcc.Program

	// An abstract of a BCC program, which provides interfaces to manage its lifetime,
	// and performance various operations like allocating input & output memory.
	wasm *wasm.Module

	// Describes the schema of the serialized output from WASM module.
	// If the data is encoded in a multi-column-format, like TLV, then the data
	// has to be decoded before writing into the data table.
	outputSchema *pg.Schema

	// The client to the database that stores Observability data.
	pgClient *pg.Client
}

// Deploy deploys eBPF+WASM module. Returns the Module object and error if failed.
func Deploy(modPB *modulepb.Module, pgClient *pg.Client) (*Module, error) {
	m := new(Module)

	m.modulePB = modPB

	ebpfProg, err := bcc.NewProgram(modPB.Ebpf)
	if err != nil {
		return nil, fmt.Errorf("while deploying, failed to create eBPF program manager, error: %v", err)
	}
	err = ebpfProg.Init()
	if err != nil {
		return nil, fmt.Errorf("while deploying, failed to initialize eBPF program manager, error: %v", err)
	}
	m.ebpf = ebpfProg

	wasmModule, err := wasm.NewWasiModule(modPB.Wasm.Code, []string{})
	if err != nil {
		return nil, fmt.Errorf("while deploying, failed to initialize WASM module, error: %v", err)
	}
	m.wasm = wasmModule
	m.outputSchema = pg.SchemaFromPB(modPB.Wasm.OutputSchema)
	m.pgClient = pgClient
	return m, nil
}

func (m *Module) StartPoll() {
	for {
		err := m.Poll()
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (m *Module) Name() string {
	return m.modulePB.Name
}

func (m *Module) Undeploy() {
	m.ebpf.Stop()
}

// Poll runs the whole process of polling data from eBPF, copying the data to WASM, reading the result from WASM.
func (m *Module) Poll() error {
	perfBufName := m.modulePB.Ebpf.PerfBufferName
	namedData := m.ebpf.Poll()
	dataItems, found := namedData[perfBufName]
	if !found {
		return fmt.Errorf("the only perf buffer '%s' is not found in polled data, %v", perfBufName, namedData)
	}

	outputDataItems := make([][]byte, 0)

	for _, data := range dataItems {
		log.Infof("lenght=%d data=%v", len(data), data)

		_, err := wasm.MallocInputBuf(m.wasm, int32(len(data)))
		if err != nil {
			return fmt.Errorf(
				"while copying polled data from eBPF to WASM, failed to malloc input buffer in WASM, error: %v",
				err,
			)
		}
		defer func() {
			err := wasm.FreeInputBuf(m.wasm)
			if err != nil {
				log.Warnf("While processing data item in WASM, failed to free input buffer, error: %v", err)
			}
		}()

		err = wasm.CopyToInputBuf(m.wasm, data)
		if err != nil {
			return fmt.Errorf(
				"while processing data in eBPF+WASM module, failed to copy data to WASM input buffer, error: %v",
				err,
			)
		}
		// The WASM function should have malloced the output buffer.
		// So here we do not malloc output buffer.
		_, err = m.wasm.Run(m.modulePB.Wasm.FnName)

		// Ensure that we free the output buffer before returning.
		// Assume the output buffer has already been allocated in the WASM function.
		defer func() {
			err := wasm.FreeOutputBuf(m.wasm)
			if err != nil {
				log.Warnf("While processing data item in WASM, failed to free output buffer, error: %v", err)
			}
		}()

		if err != nil {
			return fmt.Errorf(
				"while processing data in eBPF+WASM module, failed to run WASM function '%s', error: %v",
				m.modulePB.Wasm.FnName,
				err,
			)
		}
		// data encoding paradigm is in m.moduleDB.WasmOutputEncoding
		data, err := wasm.ReadFromOutputBuf(m.wasm)
		if err != nil {
			return fmt.Errorf("while processing data in eBPF+WASM module, failed to read output, error: %v", err)
		}
		outputDataItems = append(outputDataItems, data)
	}
	err := m.outputJSON(outputDataItems)
	if err != nil {
		return fmt.Errorf("while polling module '%s', failed to write JSON to database, error: %v", m.Name(), err)
	}
	return nil
}

func (m *Module) outputJSON(jsons [][]byte) error {
	for _, json := range jsons {
		log.Infof("length=%d json=%v", len(json), json)
		// Found in the test that there is mysterious triling null chars outputted by eBPF probes.
		// So just trim them away.
		json = bytes.Trim(json, "\x00")
		err := m.pgClient.WriteRecord([]interface{}{json}, m.outputSchema)
		if err != nil {
			return fmt.Errorf("while outputing JSON data, failed to write record to database, error: %v", err)
		}
	}
	return nil
}
