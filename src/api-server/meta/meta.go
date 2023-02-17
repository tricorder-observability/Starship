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
