package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	uuid1 := New()
	uuid2 := New()
	assert := assert.New(t)
	assert.NotEqual(uuid1, uuid2, "New() should return different results, got '%s'", uuid1)
}
