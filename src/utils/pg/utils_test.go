package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that the pgPath() returns the correct JSON path for querying JSON object.
func TestPGPath(t *testing.T) {
	assert := assert.New(t)
	idPath := []string{}
	assert.Equal("data->'metadata'->>'uid'", pgPath(idPath))

	idPath = []string{"uid"}
	assert.Equal("data->>'uid'", pgPath(idPath))

	idPath = []string{"metadata", "uid", "id"}
	assert.Equal("data->'metadata'->'uid'->>'id'", pgPath(idPath))
}
