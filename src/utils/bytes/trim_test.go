package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesTrim(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]byte("012345"), TrimAfter([]byte("012345\x00\x00"), '\x00'))
	assert.Equal([]byte("01234"), TrimAfter([]byte("012345\x00\x00"), '5'))
	assert.Equal([]byte("012345\x00\x00"), TrimAfter([]byte("012345\x00\x00"), 'B'))
	assert.Equal([]byte("012345"), TrimC([]byte("012345\x00\x00")))
}
