// Copyright (C) 2023  tricorder-observability
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
