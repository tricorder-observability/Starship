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

package grpc

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/tricorder/src/api-server/http/dao"
)

// Tests that the http service can handle request
func TestGetDeployReqForModule(t *testing.T) {
	assert := assert.New(t)
	moduleGORM := dao.ModuleGORM{ID: "test"}
	req, err := getDeployReqForModule(&moduleGORM)
	assert.Nil(err)
	assert.Equal("test", req.ModuleId)
}
