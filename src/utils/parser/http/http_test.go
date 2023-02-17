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
