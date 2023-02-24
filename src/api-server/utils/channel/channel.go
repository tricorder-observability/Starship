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

package channel

// Encloses information sent from API Server's HTTP handler to gRPC handler for more efficient job hand-off.
type DeployChannelModule struct {
	ID     string
	Status int
}

// TODO: Channels used for signalling that users have instructed API server to deploy modules.
// This should be put inside src/api-server/shared/chan.go; where HTTP server writes to this
// channel, and gRPC side waits on this channel, and triggers API server to query SQLite DB.
var (
	notifyChan chan DeployChannelModule
)

// https://www.digitalocean.com/community/tutorials/understanding-init-in-go
func init() {
	notifyChan = make(chan DeployChannelModule, 100)
}

func SendMessage(module DeployChannelModule) {
	notifyChan <- module
}

func ReceiveMessage() DeployChannelModule {
	message := <-notifyChan
	return message
}
