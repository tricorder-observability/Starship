package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen(t *testing.T) {
	assert := assert.New(t)
	req := Gen()
	assert.Equal("GET", req.Method)
}
