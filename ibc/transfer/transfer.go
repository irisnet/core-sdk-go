package transfer

import (
	"context"

	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

type transferClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return transferClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (tc transferClient) Name() string {
	return ModuleName
}

func (tc transferClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (tc transferClient) Transfer(request TransferRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	author, err := tc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgTransfer{
		SourcePort:       request.SourcePort,
		SourceChannel:    request.SourceChannel,
		Token:            request.Token,
		Sender:           author.String(),
		Receiver:         request.Receiver,
		TimeoutHeight:    request.TimeoutHeight,
		TimeoutTimestamp: request.TimeoutTimestamp,
	}
	return tc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (tc transferClient) QueryDenomTrace(request QueryDenomTraceRequest) (QueryDenomTraceResponse, sdk.Error) {
	conn, err := tc.GenConn()
	if err != nil {
		return QueryDenomTraceResponse{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).DenomTrace(
		context.Background(),
		&QueryDenomTraceRequest{
			Hash: request.Hash,
		},
	)
	if err != nil {
		return QueryDenomTraceResponse{}, sdk.Wrap(err)
	}

	return QueryDenomTraceResponse{
		DenomTrace: res.DenomTrace,
	}, nil
}

func (tc transferClient) QueryDenomTraces(request QueryDenomTracesRequest) (QueryDenomTracesResponse, sdk.Error) {
	conn, err := tc.GenConn()
	if err != nil {
		return QueryDenomTracesResponse{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).DenomTraces(
		context.Background(),
		&QueryDenomTracesRequest{
			Pagination: request.Pagination,
		},
	)
	if err != nil {
		return QueryDenomTracesResponse{}, sdk.Wrap(err)
	}

	return QueryDenomTracesResponse{
		DenomTraces: res.DenomTraces,
		Pagination:  res.Pagination,
	}, nil
}
