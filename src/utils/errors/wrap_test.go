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
	"errors"
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

// Tests that Is returns true when the error is the same.
func TestIs(t *testing.T) {
	err1 := errors.New("1")
	erra := err1
	errb := errors.New("b")

	err3 := errors.New("3")

	testCases := []struct {
		err    error
		target error
		match  bool
	}{
		{nil, nil, true},
		{err1, nil, false},
		{err1, err1, true},
		{erra, err1, true},
		{errb, err1, false},
		{err1, err3, false},
		{erra, err3, false},
		{errb, err3, false},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.match {
				t.Errorf("Is(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
			}
		})
	}
}
