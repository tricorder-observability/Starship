package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that GetURL ignores http:// prefix
func TestGetURL(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("http://localhost:8080/api/test", GetURL("localhost:8080", "/api/test"))
	assert.Equal("http://localhost:8080/api/test", GetURL("http://localhost:8080", "/api/test"))
}
