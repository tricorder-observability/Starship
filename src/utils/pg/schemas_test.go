package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that the schema created by GetJSONBTableSchema().
func TestGetTableSchema(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(&Schema{
		Name: "test_table",
		Columns: []Column{
			{
				Name: "data",
				Type: JSONB,
			},
		},
	},
		GetJSONBTableSchema("test_table"))
}
