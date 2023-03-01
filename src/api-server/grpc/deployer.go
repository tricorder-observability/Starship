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

package grpc

import (
	"encoding/json"
	"io"

	"golang.org/x/sync/errgroup"

	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/log"
	"github.com/tricorder/src/utils/sqlite"

	"github.com/tricorder/src/api-server/http/dao"
	pb "github.com/tricorder/src/api-server/pb"
	servicepb "github.com/tricorder/src/api-server/pb"
	modulepb "github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/pb/module/common"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	wasmpb "github.com/tricorder/src/pb/module/wasm"
)

// Manages the deployment of eBPF+WASM modules
type Deployer struct {
	// The DAO object that proxies with SQLite for writing and reading the serialized data.
	Module         dao.ModuleDao
	NodeAgent      dao.NodeAgentDao
	ModuleInstance dao.ModuleInstanceDao
	gLock          *lock.Lock
	waitCond       *cond.Cond
	// The list of agents connected with this Deployer.
	//
	// Each agent and this Deployer maintains a gRPC streaming channel with DeployModuleReq & DeployModuleResp
	// flow back-and-forth.
	agents []*servicepb.Agent
}

func getDeployReqForModule(module *dao.ModuleGORM) (*servicepb.DeployModuleReq, error) {
	var probeSpecs []*ebpfpb.ProbeSpec
	if len(module.EbpfProbes) > 0 {
		err := json.Unmarshal([]byte(module.EbpfProbes), &probeSpecs)
		if err != nil {
			return nil, errors.Wrap("creating DeployModuleReq for module", "unmarshal ebpf probes", err)
		}
	}

	ebpf := &ebpfpb.Program{
		Fmt:            common.Format(module.EbpfFmt),
		Lang:           common.Lang(module.EbpfLang),
		Code:           module.Ebpf,
		PerfBufferName: module.EbpfPerfBufferName,
		Probes:         probeSpecs,
	}

	var fields []*common.DataField
	if len(module.SchemaAttr) > 0 {
		err := json.Unmarshal([]byte(module.SchemaAttr), &fields)
		if err != nil {
			return nil, errors.Wrap("creatign DeployModuleReq for module", "unmarshal data fields", err)
		}
	}

	wasm := &wasmpb.Program{
		Fmt:    common.Format(module.WasmFmt),
		Lang:   common.Lang(module.WasmLang),
		FnName: module.Fn,
		OutputSchema: &common.Schema{
			Name:   module.SchemaName,
			Fields: fields,
		},
		Code: module.Wasm,
	}

	codeReq := servicepb.DeployModuleReq{
		ModuleId: module.ID,
		Module: &modulepb.Module{
			Ebpf: ebpf,
			Wasm: wasm,
		},
		Deploy: servicepb.DeployModuleReq_DEPLOY,
	}
	return &codeReq, nil
}

// DeployModule implements the only RPC of the ModuleDeployer service.
// It continuously sends deployment request to the connected agent (as client).
func (s *Deployer) DeployModule(stream servicepb.ModuleDeployer_DeployModuleServer) error {
	// The first message is sent from the client, but the remaining loops are driven by the server.
	// The server will send deploy module request for this client to work on.
	in, err := stream.Recv()
	if err == io.EOF {
		log.Warnf("Agent closed connection, this should only happens during testing; stopping ...")
		return nil
	}
	if err != nil {
		return errors.Wrap("handling DeployModule", "receive message", err)
	}

	log.Infof("Agent '%s' connected, starting module management loop ...", in.Agent.Id)
	agentNodeName := in.Agent.NodeName
	agentID := in.Agent.Id
	err = s.gLock.ExecWithLock(func() error {
		// TODO(jun): Consider return a list of nodes, and error out when there is more than 1 records
		// for this node. This is just being defensive, as ignoring that state, as the code below,
		// might result into issues that are too difficult to debug.
		node, err := s.NodeAgent.QueryByName(agentNodeName)
		if err != nil {
			// TODO(jun): Need to distinguish between query failure and getting no results.
			// It seems currently, it should return an error when there is no record for node name.
			return nil
		}
		// TODO(jun): If returning nil here means no record for this node name, then the above if statement
		// should instead return err not nil.
		if node != nil && node.State == int(pb.AgentState_ONLINE) {
			if node.AgentID == agentID {
				log.Warnf("Node '%s' agent ID '%s' was already 'ONLINE' when it connects", node.NodeName, node.AgentID)
				return nil
			}
			// There is an agent on this node with ONLINE state. And that agent is different from my ID.
			// Here we trust K8s, and assume metadata service (not yet implemented, @Daniel is working on this)
			// was slow to update the state. So we explicitly set the state to TERMINATED.
			err = s.NodeAgent.UpdateStateByName(node.NodeName, int(pb.AgentState_TERMINATED))
			if err != nil {
				return errors.Wrap("handling Agent grpc request", "update node agent state", err)
			}
			return nil
		}
		if node == nil {
			node = &dao.NodeAgentGORM{
				NodeName: agentNodeName,
				AgentID:  agentID,
				State:    int(pb.AgentState_ONLINE),
			}
			err = s.NodeAgent.SaveAgent(node)
			if err != nil {
				return errors.Wrap("handling Agent grpc request", "save new online agent", err)
			}
			return nil
		}
		// TODO(jun): The following code does not seem possible. We should log.Warnf() here to record this state.
		err = s.NodeAgent.UpdateStateByName(node.NodeName, int(pb.AgentState_ONLINE))
		if err != nil {
			return errors.Wrap("while handling Agent grpc request", "update node agent state", err)
		}
		return nil
	})

	if err != nil {
		return errors.Wrap("handling agent grpc request", "update node agent state", err)
	}

	s.agents = append(s.agents, in.Agent)

	// TODO(jun): handle the case where the node is not new, but the agent is restarted.

	var eg errgroup.Group
	// Create a goroutine to check the response from the connected agent.
	eg.Go(func() error {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				log.Warnf("Agent closed connection, this should **only** happens during testing; stopping ...")
				return nil
			}
			if err != nil {
				// If this happens, agent should re-initiate connection with API Server.
				// API Server just close the handling function and wait for reconnection.
				return errors.Wrap("handling DeployModule request", "receive mssage", err)
			}
			// TODO(yzhao): Should cache this result to an internal slice, and repeatively retry updating state.
			// The current logic will drop this state and causes redeployment of the same module.
			err = s.Module.UpdateStatusByID(result.ModuleId, int(result.State))
			if err != nil {
				log.Errorf("update code status error:%s", err.Error())
			}
		}
	})

	for {
		s.waitCond.Wait()
		undeployList, _ := s.Module.ListModuleByStatus(int(servicepb.ModuleState_DEPLOYED))
		for _, code := range undeployList {
			codeReq, err := getDeployReqForModule(&code)
			if err != nil {
				log.Fatalf("Failed to create DeployModuleReq for module ID=%s, this should not happen, "+
					"as module creation should validate module, error: %v", code.ID, err)
				return err
			}

			err = stream.Send(codeReq)
			if err != nil {
				return errors.Wrap("handling module deployment", "send message over gRCP streaming channel", err)
			}

			// TODO(yzhao): This should set the state to PENDING, or something indicating the request is sent.
			// Probably should update the IN_PROGRESS state in module_instance table.
			err = s.Module.UpdateStatusByID(code.ID, int(servicepb.ModuleState_DEPLOYED))
			if err != nil {
				// If this happens, this module's deployment will be retried next time.
				log.Errorf("Failed to update module (ID=%s) state, error: %v", code.ID, err)
			}
		}
	}
}

// NewDeployer returns a Deployer object with the input SQLite ORM client.
func NewDeployer(orm *sqlite.ORM, gLock *lock.Lock, waitCond *cond.Cond) *Deployer {
	return &Deployer{
		Module: dao.ModuleDao{
			Client: orm,
		},
		NodeAgent: dao.NodeAgentDao{
			Client: orm,
		},
		ModuleInstance: dao.ModuleInstanceDao{
			Client: orm,
		},
		waitCond: waitCond,
		gLock:    gLock,
	}
}

// RegisterDeployerService registers Deployer server instance with the gRPC fixture.
func RegisterModuleDeployerServer(f *ServerFixture, sqliteClient *sqlite.ORM, gLock *lock.Lock, waitCond *cond.Cond) {
	pb.RegisterModuleDeployerServer(f.Server, NewDeployer(sqliteClient, gLock, waitCond))
}
