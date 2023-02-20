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

package grpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tricorder/src/api-server/dao"
	pb "github.com/tricorder/src/api-server/pb"
	testutil "github.com/tricorder/src/api-server/testing"
	"github.com/tricorder/src/utils/channel"
	"github.com/tricorder/src/utils/sqlite"
)

var codeID = "9999"

// Tests that the http service can handle request
func TestService(t *testing.T) {
	testDir, _ := os.Getwd()
	testDbFilePath := testDir + "/testdata/"
	sqliteClient, _ := dao.InitSqlite(testDbFilePath)
	codeDao := dao.Module{
		Client: sqliteClient,
	}
	testutil.PrepareTricorderDBData(codeID, codeDao)
	withServerAndClient(t, sqliteClient, func(server *grpcServer, c *grpcClient) {
		in, err := c.stream.Recv()
		if err == io.EOF {
			fmt.Printf("receive stream err: %s", err.Error())
		}
		if err != nil {
			fmt.Printf("Failed to read stream from DeplyModule(), error: %v", err)
		}

		fmt.Printf("Received request to deploy module: %v", in)
		assert.Equal(t, codeID, in.ID)
		_ = os.RemoveAll(testDbFilePath)
	})
	_ = os.RemoveAll(testDir + "/tricorder.db")
}

func newDeployerServer(t *testing.T, sqliteClient *sqlite.ORM) (*grpc.Server, net.Addr) {
	lis, _ := net.Listen("tcp", ":0")
	grpcServer := grpc.NewServer()

	pb.RegisterModuleDeployerServer(grpcServer, &Deployer{
		Module: dao.Module{
			Client: sqliteClient,
		},
	})

	go func() {
		err := grpcServer.Serve(lis)
		require.NoError(t, err)
	}()

	return grpcServer, lis.Addr()
}

type grpcServer struct {
	server  *grpc.Server
	lisAddr net.Addr
}

type grpcClient struct {
	c      pb.ModuleDeployerClient
	stream pb.ModuleDeployer_DeployModuleClient
	conn   *grpc.ClientConn
}

func initializeTestServerGRPCWithOptions(t *testing.T, sqliteClient *sqlite.ORM) *grpcServer {
	server, addr := newDeployerServer(t, sqliteClient)
	return &grpcServer{
		server:  server,
		lisAddr: addr,
	}
}

func newGRPCClient(t *testing.T, addr string) *grpcClient {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	c := pb.NewModuleDeployerClient(conn)
	deployModuleStream, err := c.DeployModule(context.Background())
	if err != nil {
		log.Fatalf("Could not open stream to DeplyModule RPC at %s, %v", addr, err)
	}

	resp := pb.DeployModuleResp{
		ID: "testid",
	}
	err = deployModuleStream.Send(&resp)
	if err != nil {
		log.Fatalf("Could not send stream to DeplyModule RPC at %s, %v", addr, err)
	}

	message := channel.DeployChannelModule{
		ID:     "moduleID",
		Status: int(pb.DeploymentStatus_TO_BE_DEPLOYED),
	}
	channel.SendMessage(message)
	return &grpcClient{
		c:      c,
		stream: deployModuleStream,
		conn:   conn,
	}
}

func withServerAndClient(
	t *testing.T,
	sqliteClient *sqlite.ORM,
	actualTest func(server *grpcServer, client *grpcClient),
) {
	server := initializeTestServerGRPCWithOptions(t, sqliteClient)
	c := newGRPCClient(t, server.lisAddr.String())
	defer server.server.Stop()
	defer c.conn.Close()

	actualTest(server, c)
}
