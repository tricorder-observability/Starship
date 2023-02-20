package pb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/tricorder/src/utils/pb/testdata"
)

// Tests that protobuf message is formatted correctly.
func TestFormatOneLine(t *testing.T) {
	assert := assert.New(t)

	msg := pb.TestMessage{
		Name:    "test_name",
		Address: "test",
		Title:   "CEO",
	}

	// Golang API keeps alters spaces between fields, so have to use regexp to match.
	assert.Regexp(`name:"test_name" +address:"test" +title:"CEO"`, FormatOneLine(&msg))
	assert.Regexp(strings.Join([]string{`name: +"test_name"`, `address: +"test"`, `title: +"CEO"`}, "\n"),
		FormatMultiLine(&msg))
}
