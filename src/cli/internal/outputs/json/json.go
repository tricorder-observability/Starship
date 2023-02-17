package json

import (
	"encoding/json"
	"fmt"

	"github.com/tricorder/src/cli/internal/model"
)

func Output(data *model.Response) error {
	bytes, e := json.Marshal(data)
	if e != nil {
		return e
	}
	_, e = fmt.Printf("%v\n", string(bytes))
	return e
}
