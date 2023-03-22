// Copyright (C) 2023  Tricorder Observability
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

func TestNewWithSeparator(t *testing.T) {
	assert := assert.New(t)
	uuid1 := NewWithSeparator("===")
	assert.Regexp(`.+===.+===.+===.+===.+`, uuid1)
	assert.NotContains("-", uuid1)
}
