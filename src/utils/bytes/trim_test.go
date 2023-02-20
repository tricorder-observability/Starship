package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesTrim(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("012345", StrTrimAfter("012345\x00\x00", "\x00"))
	assert.Equal("01234", StrTrimAfter("012345\x00\x00", "5"))
	assert.Equal("0123", StrTrimAfter("012345\x00\x00", "4"))

	assert.Equal("012345", StrTrimC("012345\x00\x00"))
	assert.Equal("012345", StrTrimC("012345\x00\x00\xae"))
}
