package pb

import (
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// FormatOneLine returns an one line string as text format of the input proto message.
func FormatOneLine(m proto.Message) string {
	opts := prototext.MarshalOptions{
		Multiline: false,
	}
	return opts.Format(m)
}

func FormatMultiLine(m proto.Message) string {
	opts := prototext.MarshalOptions{
		Multiline: true,
	}
	return opts.Format(m)
}
