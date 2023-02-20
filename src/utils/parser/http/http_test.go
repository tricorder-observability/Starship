// Copyright (C) 2023  tricorder-observability
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

package http

import (
	"io"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	t.Log("Testing Parser's APIs")

	httpGetText := "GET /hello.htm HTTP/1.1\r\n" +
		"User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)\r\n" +
		"Host: www.tutorialspoint.com\r\n" +
		"Accept-Language: en-us\r\n" +
		"Accept-Encoding: gzip, deflate\r\n" +
		"Connection: Keep-Alive\r\n\r\n"

	reqs, left, err := ReadRequests([]byte(httpGetText + httpGetText))

	if err != io.EOF {
		t.Errorf("Got error: %v", err)
	}

	if left != 0 {
		t.Errorf("Expect 10 bytes left, got %d\n", left)
	}

	methods := []string{}
	for _, req := range reqs {
		methods = append(methods, req.Method)
	}

	expectedMethods := []string{"GET", "POST"}
	if reflect.DeepEqual(methods, expectedMethods) {
		t.Errorf("Expect HTTP methods %v got %v", methods, err)
		t.FailNow()
	}
}
