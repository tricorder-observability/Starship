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

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/log"
	pbutils "github.com/tricorder/src/utils/pb"

	"github.com/tricorder/src/api-server/dao"
	servicepb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/api-server/utils/channel"
	modulepb "github.com/tricorder/src/pb/module"
	"github.com/tricorder/src/pb/module/common"
	ebpfpb "github.com/tricorder/src/pb/module/ebpf"
	wasmpb "github.com/tricorder/src/pb/module/wasm"
)

// Manages the deployment of eBPF+WASM modules
type Deployer struct {
	// The DAO object that proxies with SQLite for writing and reading the serialized data.
	Module dao.Module

	// The list of agents connected with this Deployer.
	//
	// Each agent and this Deployer maintains a gRPC streaming channel with DeployModuleReq & DeployModuleResp
	// flow back-and-forth.
	agents []*servicepb.Agent
}

func getDeployReqForModule(code *dao.ModuleGORM) (*servicepb.DeployModuleReq, error) {
	var probeSpecs []*ebpfpb.ProbeSpec
	if len(code.EbpfProbes) > 0 {
		err := json.Unmarshal([]byte(code.EbpfProbes), &probeSpecs)
		if err != nil {
			return nil, errors.Wrap("creating DeployModuleReq for module", "unmarshal ebpf probes", err)
		}
	}

	ebpf := &ebpfpb.Program{
		Fmt:            common.Format(code.EbpfFmt),
		Lang:           common.Lang(code.EbpfLang),
		Code:           code.Ebpf,
		PerfBufferName: code.EbpfPerfBufferName,
		Probes:         probeSpecs,
	}

	var fields []*common.DataField
	if len(code.SchemaAttr) > 0 {
		err := json.Unmarshal([]byte(code.SchemaAttr), &fields)
		if err != nil {
			return nil, errors.Wrap("creatign DeployModuleReq for module", "unmarshal data fields", err)
		}
	}

	wasm := &wasmpb.Program{
		Fmt:    common.Format(code.WasmFmt),
		Lang:   common.Lang(code.WasmLang),
		FnName: code.Fn,
		OutputSchema: &common.Schema{
			Name:   code.SchemaName,
			Fields: fields,
		},
		Code: code.Wasm,
	}

	codeReq := servicepb.DeployModuleReq{
		ModuleId: code.ID,
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
		return nil
	}
	if err != nil {
		return errors.Wrap("handling DeployModule", "receive message", err)
	}

	log.Infof("Agent '%s' connected, starting module management loop ...", in.Agent.Id)

	s.agents = append(s.agents, in.Agent)

	var eg errgroup.Group

	// Create a goroutine to check the response from the connected agent.
	eg.Go(func() error {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				log.Warnf("Agent closed connection, this should only happens during testing; stopping ...")
				return nil
			}
			if err != nil {
				return errors.Wrap("handling DeployModule request", "receive mssage", err)
			}
			// Need to be able to correctly account which nodes are deployed, and which are
			// not deployed. TODO(yzhao & zhihui): Chat about the design.
			err = s.Module.UpdateStatusByID(result.ModuleId, int(result.State))
			if err != nil {
				log.Errorf("update code status error:%s", err.Error())
			}
		}
	})

	for {
		message := channel.ReceiveMessage()
		if message.Status != int(servicepb.DeploymentState_TO_BE_DEPLOYED) {
			continue
		}
		undeployList, _ := s.Module.ListCodeByStatus(int(servicepb.DeploymentState_TO_BE_DEPLOYED))
		for _, code := range undeployList {
			codeReq, err := getDeployReqForModule(&code)
			if err != nil {
				log.Fatalf("Failed to create DeployModuleReq for module ID=%s, this should not happen, "+
					"as module creation should validate module, error: %v", err)
				return err
			}

			err = stream.Send(codeReq)
			if err != nil {
				log.Errorf("gRPC streaming channel to agent=%s broken, error: %v", pbutils.FormatOneLine(in.Agent), err)
				return err
			}

			// TODO(yzhao): This should set the state to PENDING, or something indicating the request is sent.
			// Probably should update the IN_PROGRESS state in module_instance table.
			err = s.Module.UpdateStatusByID(code.ID, int(servicepb.DeploymentState_DEPLOYMENT_SUCCEEDED))
			if err != nil {
				log.Errorf("Failed to update module (ID=%s) state, error: %v", code.ID, err)
			}
		}
	}
}
