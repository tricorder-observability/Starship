package outputs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/cli/internal/model"
	sysutils "github.com/tricorder/src/testing/sys"
)

func TestOutPut(t *testing.T) {
	mod := model.Response{
		Code:    "200",
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
			"json output", JSON, "{\"data\":[{\"name\":\"mock-data\"}],\"code\":\"200\",\"message\":\"success\"}",
		},
		{
			"yaml output", YAML, "code: \"200\"\nmessage: success",
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
