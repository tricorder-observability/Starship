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

package meta

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/tricorder/src/utils/pg"
)

const (
	NodeTable       = "nodes"
	NameSpaceTable  = "namespaces"
	PodTable        = "pods"
	EndPointTable   = "endpoints"
	ServiceTable    = "services"
	ReplicSetTable  = "replicasets"
	DeploymentTable = "deployments"
)

var pgTables = []string{
	NodeTable,
	NameSpaceTable,
	PodTable,
	EndPointTable,
	ServiceTable,
	ReplicSetTable,
	DeploymentTable,
}

// InitResourceTables creates data tables for each and every type of resources.
func initResourceTables(pgClient *pg.Client) error {
	for _, table := range pgTables {
		if err := pgClient.CreateTable(pg.GetJSONBTableSchema(table)); err != nil {
			return fmt.Errorf("while initializing resource table '%s', failed to create table, error: %v", table, err)
		}
	}
	return nil
}

func StartWatchingResources(clientset kubernetes.Interface, pgClient *pg.Client) error {
	return NewResourceWatcher(clientset, pgClient).StartWatching()
}
