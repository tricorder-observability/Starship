package yaml

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/tricorder/src/cli/internal/model"
)

func Output(data *model.Response) error {
	bytes, e := yaml.Marshal(data)
	if e != nil {
		return e
	}
	_, e = fmt.Printf("%v", string(bytes))
	return e
}
