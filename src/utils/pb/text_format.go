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
