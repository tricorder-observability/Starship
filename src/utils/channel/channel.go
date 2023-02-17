package channel

import (
	"sync"
)

type DeployChannelModule struct {
	ID     string
	Status int
}

// TODO: Channels used for signalling that users have instructed API server to deploy modules.
// This should be put inside src/api-server/shared/chan.go; where HTTP server writes to this
// channel, and gRPC side waits on this channel, and triggers API server to query SQLite DB.
var (
	chanInstance    chan DeployChannelModule
	chanOnceManager sync.Once
)

// init chan only once
// TODO(yzhao): Can be replaced by init() in the package:
// https://www.digitalocean.com/community/tutorials/understanding-init-in-go
func initAgentChan() chan DeployChannelModule {
	chanOnceManager.Do(func() {
		chanInstance = make(chan DeployChannelModule, 100)
	})
	return chanInstance
}

func SendMessage(module DeployChannelModule) {
	initAgentChan()
	chanInstance <- module
}

func ReceiveMessage() DeployChannelModule {
	initAgentChan()
	message := <-chanInstance
	return message
}
