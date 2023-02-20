package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrTrimPrefix(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("12345", StrTrimPrefix("012345", 1))
	assert.Equal("01234", StrTrimSuffix("012345", 1))
	assert.Panics(func() { _ = StrTrimPrefix("0", 2) })
	assert.Panics(func() { _ = StrTrimSuffix("0", 2) })
}
