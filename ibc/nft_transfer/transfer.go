package nft_transfer

import (
	context "context"

	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

type nftTransferClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) nftTransferClient {
	return nftTransferClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (nft nftTransferClient) Name() string {
	return ModuleName
}

func (ntc nftTransferClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (ntc nftTransferClient) NFTTransfer(request MsgTransfer, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	author, err := ntc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	request.Sender = author.String()
	// msg := &MsgTransfer{
	// 	SourcePort:       request.SourcePort,
	// 	SourceChannel:    request.SourceChannel,
	// 	ClassId:		  request.ClassId,
	// 	TokenIds:		  request.TokenIds,
	// 	Sender:           author.String(),
	// 	Receiver:         request.Receiver,
	// 	TimeoutHeight:    request.TimeoutHeight,
	// 	TimeoutTimestamp: request.TimeoutTimestamp,
	// }
	return ntc.BuildAndSend([]sdk.Msg{&request}, baseTx)
}

func (ntc nftTransferClient) QueryClassTrace(request QueryClassTraceRequest) (QueryClassTraceResponse, sdk.Error) {
	conn, err := ntc.GenConn()
	if err != nil {
		return QueryClassTraceResponse{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).ClassTrace(context.Background(),
		&QueryClassTraceRequest{
			Hash: request.Hash,
		},
	)
	if err != nil {
		return QueryClassTraceResponse{}, sdk.Wrap(err)
	}
	return QueryClassTraceResponse{
		ClassTrace: res.ClassTrace,
	}, nil
}

func (ntc nftTransferClient) QueryClassTraces(request QueryClassTracesRequest) (QueryClassTracesResponse, sdk.Error) {
	conn, err := ntc.GenConn()
	if err != nil {
		return QueryClassTracesResponse{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).ClassTraces(context.Background(),
		&QueryClassTracesRequest{
			Pagination: request.Pagination,
		},
	)
	if err != nil {
		return QueryClassTracesResponse{}, sdk.Wrap(err)
	}
	return QueryClassTracesResponse{
		ClassTraces: res.ClassTraces,
	}, nil
}

func (nft nftTransferClient) QueryClassHash(request QueryClassHashRequest) (QueryClassHashResponse, sdk.Error) {
	conn, err := nft.GenConn()
	if err != nil {
		return QueryClassHashResponse{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).ClassHash(context.Background(),
		&QueryClassHashRequest{
			Trace: request.Trace,
		},
	)
	if err != nil {
		return QueryClassHashResponse{}, sdk.Wrap(err)
	}
	return QueryClassHashResponse{
		Hash: res.Hash,
	}, nil
}

func (nft nftTransferClient) QueryEscrowAddress(request QueryEscrowAddressRequest) (QueryEscrowAddressResponse, sdk.Error) {
	conn, err := nft.GenConn()
	if err != nil {
		return QueryEscrowAddressResponse{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).EscrowAddress(context.Background(),
		&QueryEscrowAddressRequest{
			PortId: request.PortId,
			ChannelId: request.ChannelId,
		},
	)
	if err != nil {
		return QueryEscrowAddressResponse{}, sdk.Wrap(err)
	}
	return QueryEscrowAddressResponse{
		EscrowAddress: res.EscrowAddress,
	}, nil
}

func (nft nftTransferClient) QueryParams(request QueryParamsRequest) (QueryParamsResponse, sdk.Error) {
	conn, err := nft.GenConn()
	if err != nil {
		return QueryParamsResponse{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).Params(context.Background(), &QueryParamsRequest{})
	if err != nil {
		return QueryParamsResponse{}, sdk.Wrap(err)
	}
	return QueryParamsResponse{
		Params: res.Params,
	}, nil
}

func (nft nftTransferClient) QueryPorts(request QueryPortsRequest) (QueryPortsResponse, sdk.Error) {
	conn, err := nft.GenConn()
	if err != nil {
		return QueryPortsResponse{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).Ports(context.Background(), &QueryPortsRequest{})
	if err != nil {
		return QueryPortsResponse{}, sdk.Wrap(err)
	}
	return QueryPortsResponse{
		Entries: res.Entries,
	}, nil
}
