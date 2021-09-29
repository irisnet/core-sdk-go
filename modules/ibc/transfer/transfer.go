package transfer

import (
	"context"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/codec"
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
)

type transferClient struct {
	types.BaseClient
	codec.Codec
}

func NewClient(bc types.BaseClient, cdc codec.Codec) Client {
	return transferClient{
		BaseClient: bc,
		Codec:      cdc,
	}
}

func (tc transferClient) Name() string {
	return ModuleName
}

func (tc transferClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (tc transferClient) Transfer(request TransferRequest, baseTx types.BaseTx) (ctypes.ResultTx, error) {
	author, err := tc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return ctypes.ResultTx{}, errors.Wrap(errors.ErrTodo, err.Error())
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
	return tc.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (tc transferClient) QueryDenomTrace(request QueryDenomTraceRequest) (QueryDenomTraceResponse, error) {
	conn, err := tc.GenConn()
	if err != nil {
		return QueryDenomTraceResponse{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).DenomTrace(
		context.Background(),
		&QueryDenomTraceRequest{
			Hash: request.Hash,
		},
	)
	if err != nil {
		return QueryDenomTraceResponse{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	return QueryDenomTraceResponse{
		DenomTrace: res.DenomTrace,
	}, nil
}

func (tc transferClient) QueryDenomTraces(request QueryDenomTracesRequest) (QueryDenomTracesResponse, error) {
	conn, err := tc.GenConn()
	if err != nil {
		return QueryDenomTracesResponse{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).DenomTraces(
		context.Background(),
		&QueryDenomTracesRequest{
			Pagination: request.Pagination,
		},
	)
	if err != nil {
		return QueryDenomTracesResponse{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	return QueryDenomTracesResponse{
		DenomTraces: res.DenomTraces,
		Pagination:  res.Pagination,
	}, nil
}
