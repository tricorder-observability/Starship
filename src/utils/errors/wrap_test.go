package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that Wrap returns error has the expected message.
func TestWrap(t *testing.T) {
	assert := assert.New(t)

	var err error
	wrappedErr := Wrap("while testing Wrap", "create", err)
	assert.Equal("while testing Wrap, failed to create, error: <nil>", wrappedErr.Error())
}
