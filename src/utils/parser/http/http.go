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
