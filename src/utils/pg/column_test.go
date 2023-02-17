package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that DefineColumn returns correct definition string.
func TestDefineColumn(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		column   Column
		expected string
	}{
		{
			Column{
				Name: "test",
				Type: INTEGER,
			},
			"test INTEGER",
		},
		{
			Column{
				Name:       "test",
				Type:       INTEGER,
				Constraint: PRIMARY_KEY,
			},
			"test INTEGER PRIMARY KEY",
		},
	}

	for _, c := range cases {
		got, err := DefineColumn(c.column)
		assert.Nil(err)
		assert.Equal(c.expected, got)
	}
}

// Tests that error messages of DefineColumn are as expected.
func TestDefineColumnErrors(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		c      Column
		result string
	}{
		{
			c: Column{
				Name: "test",
				Type: 100,
			},
			result: "type 'test' is not supported",
		},
		{
			c: Column{
				Name:       "test",
				Type:       INTEGER,
				Constraint: "test",
			},
			result: "constraint 'test' is not supported",
		},
	}

	for _, c := range cases {
		got, err := DefineColumn(c.c)
		assert.Equal("", got)
		assert.ErrorContains(err, "test")
	}
}
