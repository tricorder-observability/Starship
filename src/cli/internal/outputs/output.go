package outputs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tricorder/src/cli/internal/model"
	json_output "github.com/tricorder/src/cli/internal/outputs/json"
	"github.com/tricorder/src/cli/internal/outputs/table"
	yaml_output "github.com/tricorder/src/cli/internal/outputs/yaml"
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
		return json_output.Output(model)
	case YAML:
		return yaml_output.Output(model)
	case TABLE:
		return table.Output(model)
	default:
		return fmt.Errorf("unsupported output style: %s", style)
	}
}
