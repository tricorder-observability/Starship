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
