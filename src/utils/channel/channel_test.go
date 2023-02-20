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

package channel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/tricorder/src/api-server/pb"
)

func TestChannelDeployModule(t *testing.T) {
	message := DeployChannelModule{
		ID:     "moduleID",
		Status: int(pb.DeploymentStatus_TO_BE_DEPLOYED),
	}

	SendMessage(message)

	receive := ReceiveMessage()

	assert.Equal(t, true, receive.ID == message.ID)
	assert.Equal(t, true, receive.Status == message.Status)
}
