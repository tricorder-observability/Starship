package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	raw, err := os.ReadFile("./trace_return_value.c")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open BPF source file: %s\n", err)
		os.Exit(1)
	}

	req := BpfRequest{
		BpfCode:     string(raw),
		EventSize:   88,
		EventBuffer: "readline_events",
		Uprobes:     []string{"get_return_value"},
		Uretprobes: []UretprobeParam{
			UretprobeParam{
				Name:   "/bin/bash",
				Symbol: "readline",
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Marshal body failed: %v", err)
		return
	}

	// fmt.Println(body)
	url := "http://localhost:8000/bpf"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Send request failed: %v", err)
		return
	}

	fmt.Printf("Send request successful\n")
	defer resp.Body.Close()

	// Handle response
	var respBody BpfResponse
	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(rBody, &respBody); err != nil {
		panic(err)
	}

	fmt.Printf("Received response: %+v\n", respBody)
}
