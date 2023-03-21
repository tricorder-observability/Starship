package sys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortAddr(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(":100", PortAddr(100))
	assert.Equal("localhost:123", HostPortAddr("localhost", 123))
}

func TestListenTCP(t *testing.T) {
	assert := assert.New(t)
	_, addr, err := ListenTCP(0)
	assert.NoError(err)
	assert.Equal("tcp", addr.Network())
	assert.Regexp(`\[::\]:[0-9]+`, addr.String())
}
