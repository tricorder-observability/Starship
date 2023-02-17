package channel

import (
	"sync"
)

type DeployChannelModule struct {
	ID     string
	Status int
}

var (
	chanInstance    chan DeployChannelModule
	chanOnceManager sync.Once
)

// init chan only once
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
