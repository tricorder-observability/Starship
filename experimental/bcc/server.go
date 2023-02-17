package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	bpf "github.com/iovisor/gobpf/bcc"
)

// AttachUretprobe attaches a uretprobe fd to the symbol in the library or binary 'name'
// The 'name' argument can be given as either a full library path (/usr/lib/..),
// a library without the lib prefix, or as a binary with full path (/bin/bash)
// A pid can be given to attach to, or -1 to attach to all processes
//
// func (bpf *Module) AttachUretprobe(name, symbol string, fd, pid int)
type UretprobeParam struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Pid    int32  `json:"pid"`
}

type BpfRequest struct {
	BpfCode     string `json:"bpf_code"`
	EventSize   int    `json:"event_size"`
	EventBuffer string `json:"event_buffer"`

	Kprobes     []string         `json:"kprobes"`
	Kretprobes  []string         `json:"kretprobes"`
	Uprobes     []string         `json:"uprobes"`
	Uretprobes  []UretprobeParam `json:"uretprobes"`
	Tracepoints []string         `json:"tracepoints"`
}

type BpfResponse struct {
	Message string `json:"message"`
}

func serveAndHandle() {
	http.HandleFunc("/bpf", handleBPFRequest)
	// http.HandleFunc("/headers", headers)

	addr := "0.0.0.0:8000"

	fmt.Printf("Serving on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

func handleBPFRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("hello\n")
	var req BpfRequest

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(reqBody, &req); err != nil {
		panic(err)
	}

	fmt.Printf("Received request: %+v\n", req)
	respBody := BpfResponse{Message: ""}

	// Request validation
	if len(req.BpfCode) == 0 {
		respBody.Message = "Invalid BPF code"
		resp, err := json.Marshal(&respBody)
		if err != nil {
			panic(err)
		}

		w.Write(resp)
		return
	}

	respBody.Message = "Request accepted"
	resp, err := json.Marshal(&respBody)
	if err != nil {
		panic(err)
	}

	w.Write(resp)

	go processingBPF(&req)
}

func processingBPF(req *BpfRequest) {
	// Create new module with the raw BPF (C) code in the request
	m := bpf.NewModule(req.BpfCode, []string{})
	defer m.Close()

	// Attach uprobes/uretprobes
	if len(req.Uprobes) > 0 {
		if err := attachUserProvidedUprobes(req, m); err != nil {
			panic(err)
		}
	}

	// Read events
	fmt.Printf("Reading events from buffer %s\n", req.EventBuffer)
	table := bpf.NewTable(m.TableId(req.EventBuffer), m)

	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fmt.Printf("%10s\t%s\n", "PID", "COMMAND")
	go func() {
		type readlineEvent struct {
			Pid uint32
			Str [84]byte
		}

		var event readlineEvent
		for {
			data := <-channel
			err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &event)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}
			// Convert C string (null-terminated) to Go string
			comm := string(event.Str[:bytes.IndexByte(event.Str[:], 0)])
			fmt.Printf("%10d\t%s\n", event.Pid, comm)
		}
	}()

	perfMap.Start()
	<-sig
	perfMap.Stop()
}

func attachUserProvidedUprobes(req *BpfRequest, m *bpf.Module) error {
	tmpFD := -1 // just for test

	for _, fnName := range req.Uprobes {
		fmt.Printf("Loading uprobe function from user provided BPF: %s\n", fnName)

		fd, err := m.LoadUprobe(fnName)
		if err != nil {
			return fmt.Errorf("failed to load uprobe function from user provided BPF: %s\n", err)
		}

		tmpFD = fd
	}

	for _, param := range req.Uretprobes {
		name := param.Name
		symbol := param.Symbol
		pid := -1 // TODO: not used yet

		fmt.Printf("Attaching uretprobe from user provided BPF to %s:%s\n", name, symbol)
		if err := m.AttachUretprobe(name, symbol, tmpFD, pid); err != nil {
			return fmt.Errorf("failed to attach return_value: %s\n", err)
		}
	}

	return nil
}
