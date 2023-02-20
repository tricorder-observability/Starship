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

package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesTrim(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]byte("012345"), TrimAfter([]byte("012345\x00\x00"), '\x00'))
	assert.Equal([]byte("01234"), TrimAfter([]byte("012345\x00\x00"), '5'))
	assert.Equal([]byte("012345\x00\x00"), TrimAfter([]byte("012345\x00\x00"), 'B'))
	assert.Equal([]byte("012345"), TrimC([]byte("012345\x00\x00")))
}
