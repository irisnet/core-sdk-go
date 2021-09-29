package transfer

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/types"
)

// expose transfer module api for user
type Client interface {
	types.Module

	Transfer(request TransferRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)

	QueryDenomTrace(request QueryDenomTraceRequest) (QueryDenomTraceResponse, error)
	QueryDenomTraces(request QueryDenomTracesRequest) (QueryDenomTracesResponse, error)
}

type TransferRequest struct {
	SourcePort       string     `json:"source_port"`
	SourceChannel    string     `json:"source_channel"`
	Token            types.Coin `json:"token"`
	Sender           string     `json:"sender"`
	Receiver         string     `json:"receiver"`
	TimeoutHeight    Height     `json:"timeout_height"`
	TimeoutTimestamp uint64     `json:"timeout_timestamp"`
}
