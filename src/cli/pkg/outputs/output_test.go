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

package outputs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/cli/pkg/model"
	sysutils "github.com/tricorder/src/testing/sys"
)

func TestOutPut(t *testing.T) {
	mod := model.Response{
		Code:    200,
		Message: "success",
		Data: []map[string]interface{}{
			{
				"name": "mock-data",
			},
		},
	}

	resp, err := json.Marshal(mod)
	assert.Nil(t, err)

	cases := []struct {
		caseStr     string
		outPutStyle string
		expected    string
	}{
		{
			"json output", JSON, "{\"data\":[{\"name\":\"mock-data\"}],\"code\":200,\"message\":\"success\"}",
		},
		{
			"yaml output", YAML, "code: 200\nmessage: success",
		},
		{
			"table output", TABLE, "+-----------+\n|   NAME    |\n+-----------+\n| mock-data |\n+-----------+\n",
		},
	}
	assert := assert.New(t)
	for _, sc := range cases {
		out := sysutils.CaptureStdout(func() {
			err := Output(sc.outPutStyle, resp)
			assert.Nil(err)
		})
		assert.Contains(out, sc.expected)
	}
}
