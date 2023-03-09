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

package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/cli/pkg/model"
	sysutils "github.com/tricorder/src/testing/sys"
)

func TestOutput(t *testing.T) {
	m := &model.Response{
		Code:    "200",
		Message: "success",
		Data: []map[string]interface{}{
			{
				"name": "mock-data",
			},
		},
	}

	assert := assert.New(t)
	out := sysutils.CaptureStdout(func() {
		err := Output(m)
		assert.Nil(err)
	})
	assert.Contains(out, "code: \"200\"\nmessage: success")
}
