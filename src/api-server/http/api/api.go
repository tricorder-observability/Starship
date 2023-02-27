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

package api

// API path components.
const (
	ROOT            = "/api"
	LIST_MODULE     = "/listModule"
	CREATE_MODULE   = "/createModule"
	DEPLOY_MODULE   = "/deployModule"
	UNDEPLOY_MODULE = "/undeployModule"
	DELETE_MODULE   = "/deleteModule"

	LIST_MODULE_PATH     = ROOT + LIST_MODULE
	CREATE_MODULE_PATH   = ROOT + CREATE_MODULE
	DEPLOY_MODULE_PATH   = ROOT + DEPLOY_MODULE
	UNDEPLOY_MODULE_PATH = ROOT + UNDEPLOY_MODULE
	DELETE_MODULE_PATH   = ROOT + DELETE_MODULE
)

// GetURL returns a http URL that corresponds to the requested path.
// The path has to start with '/'
func GetURL(addr, path string) string {
	return "http://" + addr + path
}
