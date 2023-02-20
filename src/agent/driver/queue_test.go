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

package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueue tests Enqueue and Dequeue methods
func TestQueue(t *testing.T) {
	assert := assert.New(t)

	q := NewQueue(10)
	data, err := q.Dequeue()
	assert.Nil(data)
	assert.NotNil(err)

	err = q.Enqueue([]byte("01234"))
	assert.Nil(err)

	err = q.Enqueue([]byte("56789"))
	assert.Nil(err)

	// Can write empty data
	err = q.Enqueue([]byte{})
	assert.Nil(err)

	err = q.Enqueue([]byte("01234"))
	assert.NotNil(err)

	data, err = q.Dequeue()
	assert.Nil(err)
	assert.Equal(data, []byte("01234"))

	data, err = q.Dequeue()
	assert.Nil(err)
	assert.Equal(data, []byte("56789"))

	data, err = q.Dequeue()
	assert.Nil(err)
	assert.Equal(data, []byte{})
}
