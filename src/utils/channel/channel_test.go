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
