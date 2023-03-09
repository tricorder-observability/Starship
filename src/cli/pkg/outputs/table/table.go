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

package table

import (
	"encoding/json"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/tricorder/src/cli/pkg/model"
)

func Output(resp *model.Response) error {
	var stringMapArrays []map[string]string

	bytes, _ := json.Marshal(resp.Data)
	_ = json.Unmarshal(bytes, &stringMapArrays)

	if len(stringMapArrays) < 1 {
		return nil
	}

	var header []string

	for k := range stringMapArrays[0] {
		header = append(header, k)
	}

	var dataT [][]string

	for _, objMap := range stringMapArrays {
		var datum []string
		for _, key := range header {
			datum = append(datum, objMap[key])
		}
		dataT = append(dataT, datum)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(dataT)
	table.Render()

	return nil
}
