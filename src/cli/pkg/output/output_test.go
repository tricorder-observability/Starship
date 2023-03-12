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

package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/api-server/http"
	"github.com/tricorder/src/api-server/http/dao"
	"github.com/tricorder/src/testing/sys"
)

func TestOutPut(t *testing.T) {
	mod := http.ListModuleResp{
		http.HTTPResp{
			Code:    200,
			Message: "success",
		},
		[]dao.ModuleGORM{
			{
				Name: "mock-data",
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
			"json output", JSON, `{"code":200,"message":"success","data":[{"name":"mock-data"}]}`,
		},
		{
			"yaml output", YAML, "code: 200\n  message: success",
		},
		{
			"table output", TABLE, "+-----------+\n|   NAME    |\n+-----------+\n| mock-data |\n+-----------+\n",
		},
	}
	assert := assert.New(t)
	for _, sc := range cases {
		out := sys.CaptureStdout(func() {
			err := Print(sc.outPutStyle, resp)
			assert.Nil(err)
		})
		assert.Contains(out, sc.expected)
	}
}
