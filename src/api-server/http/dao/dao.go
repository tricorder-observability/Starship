package dao

import "github.com/tricorder/src/utils/sqlite"

// Dao stores the objects to access various data for different types of object in SQLite DB.
type Dao struct {
	// Stores the eBPF+WASM modules, which describes the eBPF source code, WASM source and binary code.
	Module ModuleDao
	// Stores the agents running as daemonset on every node.
	NodeAgent NodeAgentDao
	// Stores the module instances should be deployed on each and every agent.
	ModuleInstance ModuleInstanceDao
}

// NewDao returns the Dao object for accessing the data.
func NewDao(sqliteClient *sqlite.ORM) Dao {
	return Dao{
		Module:         ModuleDao{Client: sqliteClient},
		NodeAgent:      NodeAgentDao{Client: sqliteClient},
		ModuleInstance: ModuleInstanceDao{Client: sqliteClient},
	}
}
