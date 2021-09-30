package client

import (
	"github.com/prometheus/common/promlog"
	"google.golang.org/grpc"

	"github.com/irisnet/core-sdk-go/types"
)

type grpcClient struct {
	clientConn *grpc.ClientConn
}

func NewGRPCClient(url string) types.GRPCClient {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	clientConn, err := grpc.Dial(url, dialOpts...)
	if err != nil {
		_ = promlog.New(&promlog.Config{}).Log(err.Error())
		panic(err)
	}
	return &grpcClient{clientConn: clientConn}
}

func (g grpcClient) GenConn() (*grpc.ClientConn, error) {
	return g.clientConn, nil
}
