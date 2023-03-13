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

package lock

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecWithLock(t *testing.T) {
	assert := assert.New(t)
	i := 5

	lock := NewLock()
	fn := func() error {
		if i == 1 {
			return nil
		}
		return fmt.Errorf("i not equal 1")
	}

	assert.NotNil(lock.ExecWithLock(fn), "fn must return error")

	fn = func() error {
		if i == 5 {
			return nil
		}
		return fmt.Errorf("i not equal 5")
	}
	assert.Nil(lock.ExecWithLock(fn), "fn must return nil")

	// TODO(jun): add test for synconized access
}
