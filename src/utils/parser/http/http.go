// Package http provides basic parsing API for HTTP network traffic.
package http

import (
	"bufio"
	"bytes"
	"net/http"
)

// ReadRequests reads HTTP request from the input bytes, returns a list of requests, remaining byte count,
// and error if failed.
func ReadRequests(data []byte) ([]*http.Request, int, error) {
	var reqs []*http.Request

	bytesReader := bytes.NewReader(data)
	bufioReader := bufio.NewReader(bytesReader)

	req, err := http.ReadRequest(bufioReader)
	// Keep tracking how many bytes left, after the loop terminates, the value before the last loop iteration will be
	// returned.
	left := bytesReader.Len() + bufioReader.Buffered()
	for ; err == nil; req, err = http.ReadRequest(bufioReader) {
		reqs = append(reqs, req)
		left = bytesReader.Len() + bufioReader.Buffered()
	}

	return reqs, left, err
}
