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
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"

	"github.com/tricorder/src/api-server/http"
)

const (
	JSON  = "json"
	YAML  = "yaml"
	TABLE = "table"
)

func printJSON(data *http.ListModuleResp) error {
	bytes, e := json.Marshal(data)
	if e != nil {
		return e
	}
	_, e = fmt.Printf("%v\n", string(bytes))
	return e
}

func printTable(resp *http.ListModuleResp) error {
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

func printYAML(data *http.ListModuleResp) error {
	bytes, e := yaml.Marshal(data)
	if e != nil {
		return e
	}
	_, e = fmt.Printf("%v", string(bytes))
	return e
}

// Print writes output to the console.
func Print(style string, resp []byte) error {
	var model *http.ListModuleResp
	err := json.Unmarshal(resp, &model)
	if err != nil {
		return err
	}
	if len(style) == 0 {
		style = YAML
	}
	switch strings.ToLower(style) {
	case JSON:
		return printJSON(model)
	case YAML:
		return printYAML(model)
	case TABLE:
		return printTable(model)
	default:
		return fmt.Errorf("unsupported output style: %s", style)
	}
}
