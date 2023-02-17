package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"

	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/testing/pg"
	upg "github.com/tricorder/src/utils/pg"
)

// TestProcessCollectorIntegration scenario
// Firstly create a K8s with an example pod having container; then start ProcessCollector gRPC server side and client
// side, interact with each other, finally a process info should be expected to insert into DB.
func TestProcessCollectorIntegration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	clientset := fake.NewSimpleClientset()

	// Ensure there is a pod in K8s and with necessary info
	pod1 := &corev1.Pod{}
	pod1.Name = "pod1"
	// Set this pod node name is 'node1'
	pod1.Spec.NodeName = "node1"
	pod1.UID = types.UID("podUID1")
	// Set this pod status is running
	ps := corev1.PodStatus{Phase: corev1.PodRunning}
	// Set container name is 'containerName1' and ID is "123"
	cs := corev1.ContainerStatus{Name: "containerName1", ContainerID: "123"}
	ps.ContainerStatuses = []corev1.ContainerStatus{cs}
	pod1.Status = ps
	pod1.UID = types.UID("uid1")
	// Put this pod into K8s
	_, err := clientset.CoreV1().Pods(corev1.NamespaceDefault).Create(context.TODO(), pod1, metav1.CreateOptions{})
	assert.Nil(err)

	cleaner, pgClient, err := pg.LaunchContainer()
	require.Nil(err)
	defer func() {
		assert.Nil(cleaner())
	}()

	// Start ProcessCollector server side in normal flow
	grpcLis, err := net.Listen("tcp", ":50051")
	assert.Nil(err)
	defer grpcLis.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterProcessCollectorServer(grpcServer, NewPIDCollector(clientset, pgClient))
	go func() {
		err = grpcServer.Serve(grpcLis)
	}()

	// Mock a ProcessCollector client side interact with server side
	grpcConn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(err)

	processCollectorClient := pb.NewProcessCollectorClient(grpcConn)
	clientStream, err := processCollectorClient.ReportProcess(context.Background())
	assert.Nil(err)

	// Send node name to server side
	node := &pb.ProcessWrapper_NodeName{NodeName: pod1.Spec.NodeName}
	err = clientStream.Send(&pb.ProcessWrapper{Msg: node})
	assert.Nil(err)

	// Ensure server side handle this request and push a feedback
	time.Sleep(1 * time.Second)
	receivedContainerInfo, err := clientStream.Recv()
	assert.Nil(err)
	assert.Equal(pod1.Status.ContainerStatuses[0].ContainerID, receivedContainerInfo.Id)

	// Client side mock its procList by this containerInfo
	pi := &pb.Process{Id: 123456}
	process := &pb.ProcessWrapper_Process{Process: &pb.ProcessInfo{
		ProcList:  []*pb.Process{pi},
		Container: &pb.ContainerInfo{Id: cs.ContainerID},
	}}

	// Client side send processInfo to server side
	err = clientStream.Send(&pb.ProcessWrapper{Msg: process})
	assert.Nil(err)

	// Ensure msg received in server side and write to DB
	time.Sleep(1 * time.Second)
	resultInDB := []*pb.ProcessInfo{}
	err = pgClient.JSON().List(procInfoTableName, &resultInDB)
	assert.Nil(err)
	assert.Equal(1, len(resultInDB))

	// ContainerID in DB should be pod1's containerID
	assert.Equal(pod1.Status.ContainerStatuses[0].ContainerID, resultInDB[0].Container.Id)
}

// Test processInfo table's UUID works based on idPath of Postgres json
func TestIdPath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cleaner, pgClient, err := pg.LaunchContainer()
	require.Nil(err)
	defer func() {
		assert.Nil(cleaner())
	}()

	err = pgClient.CreateTable(upg.GetJSONBTableSchema(procInfoTableName))
	assert.Nil(err)

	id := "abcdefg"
	pi := &pb.ProcessInfo{Container: &pb.ContainerInfo{Id: id}}
	pi.Container.Name = "@#$%^&*()_+|"
	value, _ := json.Marshal(pi)
	err = pgClient.JSON().Upsert(procInfoTableName, id, value, idPath...)
	assert.Nil(err)

	result1 := pb.ProcessInfo{}
	err = pgClient.JSON().Get(procInfoTableName, &result1, fmt.Sprintf("WHERE data->'container'->>'id'='%s'", id))
	assert.Nil(err)
	// Check result upserted
	assert.Equal(pi.Container.Name, result1.Container.Name)
}
