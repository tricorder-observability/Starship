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
