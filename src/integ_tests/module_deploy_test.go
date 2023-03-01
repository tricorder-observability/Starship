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

package integ_tests

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tricorder/src/testing/bazel"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/lock"
	"github.com/tricorder/src/utils/log"

	apiservergrpc "github.com/tricorder/src/api-server/grpc"
	"github.com/tricorder/src/api-server/http/dao"
	pb "github.com/tricorder/src/api-server/pb"
	testutil "github.com/tricorder/src/api-server/testing"
)

var moduleID = "9999"

// Tests that the http service can handle request
func TestService(t *testing.T) {
	assert := assert.New(t)

	testDir := bazel.CreateTmpDir()
	sqliteClient, err := dao.InitSqlite(testDir)
	assert.Nil(err)
	testutil.PrepareTricorderDBData(moduleID, dao.ModuleDao{Client: sqliteClient})

	gLock := lock.NewLock()
	waitCond := cond.NewCond()

	f, err := apiservergrpc.NewServerFixture(0)
	if err != nil {
		log.Fatalf("Failed to create gRPC server fixture on :0")
	}

	RegisterModuleDeployerServer(f, sqliteClient, gLock, waitCond)
	go func() {
		err := f.Serve()
		if err != nil {
			log.Fatalf("Failed to start gRPC server, error: %v", err)
		}
	}()

	c := newGRPCClient(f.addr.String())
	defer f.Server.Stop()
	defer c.conn.Close()

	// This is completely broken, but needed to unblock API Server's internal conditional waiting.
	// Ideally we should send a request to API Server and let API Server's internal logic triggers conditional variable's
	// broadcasting.
	waitCond.Broadcast()

	in, err := c.stream.Recv()
	if err == io.EOF {
		fmt.Printf("receive stream err: %s", err.Error())
	}
	if err != nil {
		fmt.Printf("Failed to read stream from DeplyModule(), error: %v", err)
	}

	fmt.Printf("Received request to deploy module: %v", in)
	assert.Equal(moduleID, in.ModuleId)
}

type deployerClient struct {
	client pb.ModuleDeployerClient
	stream pb.ModuleDeployer_DeployModuleClient
	conn   *grpc.ClientConn
}

func newGRPCClient(addr string) *deployerClient {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect gRPC at '%s', error: %v", addr, err)
	}

	c := pb.NewModuleDeployerClient(conn)
	deployModuleStream, err := c.DeployModule(context.Background())
	if err != nil {
		log.Fatalf("Could not open stream to DeplyModule RPC at %s, %v", addr, err)
	}

	resp := pb.DeployModuleResp{
		ModuleId: "testid",
		Agent:    &pb.Agent{Id: "agent"},
	}
	err = deployModuleStream.Send(&resp)
	if err != nil {
		log.Fatalf("Could not send stream to DeplyModule RPC at %s, %v", addr, err)
	}

	return &deployerClient{
		client: c,
		stream: deployModuleStream,
		conn:   conn,
	}
}
