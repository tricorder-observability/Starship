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

// Package http provides API to generate random HTTP requests
package http

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/tricorder/src/utils/common"
)

// Gen returns a http request with a few fields randomly generated.
func Gen() http.Request {
	urlstr := "http://localhost:8080/" + common.RandStr(10)
	bodystr := common.RandStr(20)
	stringReader := strings.NewReader(bodystr)
	stringReadCloser := io.NopCloser(stringReader)
	urlResult, _ := url.Parse(urlstr)
	return http.Request{
		Method: "GET",
		URL:    urlResult,
		Proto:  "HTTP/1.1",
		Header: http.Header{
			"Accept-Encoding": {"gzip, deflate"},
			"Accept-Language": {"en-us"},
			"Request-Id":      {uuid.New().String()},
		},
		Body: stringReadCloser,
	}
}
