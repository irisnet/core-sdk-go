package client

import (
	grpc1 "github.com/gogo/protobuf/grpc"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

type grpcClient struct {
	clientConn *grpc1.ClientConn
}

func NewGRPCClient(url string) types.GRPCClient {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	clientConn, err := grpc.Dial(url, dialOpts...)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	conn := grpc1.ClientConn(clientConn)

	return &grpcClient{clientConn: &conn}
}

func (g grpcClient) GenConn() (*grpc1.ClientConn, error) {
	return g.clientConn, nil
}
