package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrTrim(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("12345", StrTrimPrefix("012345", 1))
	assert.Equal("01234", StrTrimSuffix("012345", 1))
	assert.Panics(func() { _ = StrTrimPrefix("0", 2) })
	assert.Panics(func() { _ = StrTrimSuffix("0", 2) })

	assert.Equal("012345", StrTrimAfter("012345\x00\x00", "\x00"))
	assert.Equal("01234", StrTrimAfter("012345\x00\x00", "5"))
	assert.Equal("0123", StrTrimAfter("012345\x00\x00", "4"))

	assert.Equal("012345", StrTrimC("012345\x00\x00"))
	assert.Equal("012345", StrTrimC("012345\x00\x00\xae"))
}
