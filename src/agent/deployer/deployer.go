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

// Package deployer implements the agent's logic for connecting to the API Server's ModuleDeployer service.
package deployer

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tricorder/src/utils/grpcerr"
	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/agent/driver"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/uuid"

	pb "github.com/tricorder/src/api-server/pb"
)

// Deployer manages the communication with API Server:
// * Receive instructions to deploy modules
// * Reply deployment status
type Deployer struct {
	uuid string

	// The remote API server's address. This deployer receives deployment request from the API server.
	apiServerAddr string

	// Key is the eBPF+WASM module's ID, value is the Module object.
	// The Module object keeps track of the module's deployment state.
	idDeployMap map[string]*driver.Module

	grpcConn *grpc.ClientConn
	client   pb.ModuleDeployerClient
	stream   pb.ModuleDeployer_DeployModuleClient

	// The client to the database instance, which is for the eBPF+WASM module to write data.
	PGClient *pg.Client
}

func (s *Deployer) ConnectToAPIServer(apiServerAddr string) error {
	log.Infof("connecting to API server at %s", apiServerAddr)

	s.uuid = uuid.New()
	s.apiServerAddr = apiServerAddr
	s.idDeployMap = make(map[string]*driver.Module)
	grpcConn, err := grpc.Dial(apiServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to API server at '%s', error: %v", apiServerAddr, err)
	}
	s.grpcConn = grpcConn
	client := pb.NewModuleDeployerClient(grpcConn)
	s.client = client
	return nil
}

// InitModuleDeployLink connects with the module module deployer's stream gRPC service.
// And sends the first message to the server to inform the server about its own identity.
func (s *Deployer) InitModuleDeployLink() error {
	log.Infof("initializing stream connection with ModuleDeployer at %s", s.apiServerAddr)

	deployModuleStream, err := s.client.DeployModule(context.Background())
	if err != nil {
		return fmt.Errorf("could not open stream to DeplyModule RPC at %s, %v", s.apiServerAddr, err)
	}
	s.stream = deployModuleStream

	resp := pb.DeployModuleResp{
		ID: s.uuid,
	}

	err = s.stream.Send(&resp)
	if err != nil {
		return err
	}
	return nil
}

// StartModuleDeployLoop continuously polling server
// The gRPC streaming channel should always be working, otherwise, agent just crash and restart.
// TODO(yzhao): We need to implement a graceful reconnection to ensure data remains available during the time when api
// server is unavailable, could happen when api server is being restarted.
func (s *Deployer) StartModuleDeployLoop() error {
	var eg errgroup.Group
	eg.Go(func() error {
		for {
			in, err := s.stream.Recv()
			if err != nil {
				return fmt.Errorf("failed to read stream from DeplyModule(), error: %v", err)
			}

			log.Infof("received request to deploy module. ID: [%s], Name: [%s]", in.ID, in.Name)
			log.Debugf("received request to deploy module: %v", in)

			if in.Deploy == pb.DeployModuleReq_DEPLOY {
				err := s.deployModule(in)
				resp := createDeployModuleResp(in.ID, err)
				err = s.sendResp(resp)
				if grpcerr.IsUnavailable(err) {
					log.Fatalf("streaming connection with api-server is broken, error: %v", err)
				}
			} else if in.Deploy == pb.DeployModuleReq_UNDEPLOY {
				err := s.undeployModlue(in)
				resp := createDeployModuleResp(in.ID, err)
				err = s.sendResp(resp)
				if grpcerr.IsUnavailable(err) {
					log.Fatalf("streaming connection with api-server is broken, error: %v", err)
				}
			}
		}
	})
	return eg.Wait()
}

func (s *Deployer) Stop() {
	err := s.stream.CloseSend()
	if err != nil {
		log.Errorf("failed to Close stream, error: %v", err)
	}
	s.grpcConn.Close()
}

func (s *Deployer) deployModule(in *pb.DeployModuleReq) error {
	if _, found := s.idDeployMap[in.ID]; found {
		log.Warnf("Module '%s' was already deployed, skip ...", in.ID)
		// TODO(yzhao): Might consider returning an error value to distinguish from other errors.
		return nil
	}
	// deployer create a deployment and driver will start this deploys logical
	deployment, err := driver.Deploy(in.Module, s.PGClient)
	if err != nil {
		return fmt.Errorf("while deploying module '%s', failed to deploy, error: %v", in.ID, err)
	}
	s.idDeployMap[in.ID] = deployment

	// This will start a loop to continuously polling perf buffer and feeding data to WASM.
	// And then write them into database.
	go deployment.StartPoll()
	return nil
}

func (s *Deployer) undeployModlue(in *pb.DeployModuleReq) error {
	d, ok := s.idDeployMap[in.ID]
	if !ok {
		return fmt.Errorf("while undeploying module ID '%s', could not find deployment record", in.ID)
	}

	log.Infof("Prepare undeploy module [ID: %s], [Name: %s]", in.ID, d.Name())

	d.Undeploy()
	delete(s.idDeployMap, in.ID)
	return nil
}

// createDeployModuleResp returns a response message to describe the results of a module deployment operation.
func createDeployModuleResp(id string, err error) *pb.DeployModuleResp {
	resp := pb.DeployModuleResp{
		ID:     id,
		Status: pb.DeploymentStatus_DEPLOYMENT_SUCCEEDED,
	}
	if err != nil {
		resp.Status = pb.DeploymentStatus_DEPLOYMENT_FAILED
		resp.Desc = err.Error()
	}
	return &resp
}

func (s *Deployer) sendResp(resp *pb.DeployModuleResp) error {
	// Preserve the original error message, as it's needed to check the status code.
	return s.stream.Send(resp)
}
