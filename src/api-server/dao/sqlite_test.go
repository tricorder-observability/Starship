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

package dao

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tricorder/src/testing/bazel"
)

func TestInitSqlitGorm(t *testing.T) {
	require := require.New(t)

	testDbFilePath := bazel.GetTmpFile()
	_, err := InitSqlite(testDbFilePath)
	require.Nil(err)
}
