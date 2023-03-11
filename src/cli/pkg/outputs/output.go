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
	"fmt"
	"strings"

	"github.com/tricorder/src/cli/pkg/model"
	jsonoutput "github.com/tricorder/src/cli/pkg/outputs/json"
	"github.com/tricorder/src/cli/pkg/outputs/table"
	"github.com/tricorder/src/cli/pkg/outputs/yaml"
)

const (
	JSON  = "json"
	YAML  = "yaml"
	TABLE = "table"
)

// Output writes output to the console.
func Output(style string, resp []byte) error {
	var model *model.Response
	err := json.Unmarshal(resp, &model)
	if err != nil {
		return err
	}
	if len(style) == 0 {
		style = YAML
	}
	switch strings.ToLower(style) {
	case JSON:
		return jsonoutput.Output(model)
	case YAML:
		return yaml.Output(model)
	case TABLE:
		return table.Output(model)
	default:
		return fmt.Errorf("unsupported output style: %s", style)
	}
}
