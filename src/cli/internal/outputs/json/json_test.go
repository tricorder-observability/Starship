package json

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tricorder/src/cli/internal/model"
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
	assert.Contains(out, "{\"data\":[{\"name\":\"mock-data\"}],\"code\":\"200\",\"message\":\"success\"}")
}
