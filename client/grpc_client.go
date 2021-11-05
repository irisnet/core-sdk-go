package client

import (
	grpc1 "github.com/gogo/protobuf/grpc"

	"github.com/prometheus/common/promlog"
	"google.golang.org/grpc"

	"github.com/irisnet/core-sdk-go/types"
)

type grpcClient struct {
	clientConn grpc1.ClientConn
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
	conn := grpc1.ClientConn(clientConn)
	return grpcClient{clientConn: conn}
}

func (g grpcClient) GenConn() (grpc1.ClientConn, error) {
	return g.clientConn, nil
}
