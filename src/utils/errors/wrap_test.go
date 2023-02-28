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

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that Wrap returns error has the expected message.
func TestWrap(t *testing.T) {
	assert := assert.New(t)

	var err error
	wrappedErr := Wrap("testing Wrap", "create", err)
	assert.Equal("while testing Wrap, failed to create, error: <nil>", wrappedErr.Error())
}

// Tests that New returns error has the expected message.
func TestNew(t *testing.T) {
	assert := assert.New(t)

	wrappedErr := New("while testing Wrap", "create")
	assert.Equal("while testing Wrap, failed to create", wrappedErr.Error())
}
