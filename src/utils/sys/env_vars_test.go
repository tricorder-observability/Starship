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

package sys

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests that EnvVars() returned env vars.
func TestEnvVar(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("FOO", "1")
	envVars := EnvVars()
	fooVal, found := envVars["FOO"]
	assert.True(found)
	assert.Equal("1", fooVal)
}
