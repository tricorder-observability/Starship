// Copyright (C) 2023  tricorder-observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
