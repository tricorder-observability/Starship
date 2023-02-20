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
