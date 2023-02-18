package pb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/tricorder/src/utils/pb/testdata"
)

// Tests that json_name annotation's effects.
func TestJSONName(t *testing.T) {
	assert := assert.New(t)

	msg := pb.TestMessage{
		Name: "test_name",
	}

	jsonText := "{\n" +
		`  "test_json_name": +"test_name"` + "\n" +
		"}"
	assert.Regexp(jsonText, protojson.Format(&msg))

	json1 := `{"test_json_name": "test_name"}`
	msg1 := pb.TestMessage{}
	assert.Nil(protojson.Unmarshal([]byte(json1), &msg1))
	assert.Equal(&msg, &msg1)

	// Default json name is exactly the same as the protobuf field name
	json2 := `{"name": "test_name"}`
	msg2 := pb.TestMessage{}
	assert.Nil(protojson.Unmarshal([]byte(json2), &msg2))
	assert.Equal(&msg, &msg2)
}
