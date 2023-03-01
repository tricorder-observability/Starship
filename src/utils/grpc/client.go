package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tricorder/src/utils/errors"
)

// DialInsecure returns a gRPC connection and error if failed.
func DialInsecure(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap("dialing server at "+addr, "dial without credentials", err)
	}
	return conn, nil
}
