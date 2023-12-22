package nft_transfer

import (
	sdk "github.com/irisnet/core-sdk-go/types"
)

// expose transfer module api for user
type Client interface {
	sdk.Module
	NFTTransfer(request MsgTransfer, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	QueryPorts(request QueryPortsRequest) (QueryPortsResponse, sdk.Error)
	QueryParams(request QueryParamsRequest) (QueryParamsResponse, sdk.Error)
	QueryClassHash(request QueryClassHashRequest) (QueryClassHashResponse, sdk.Error)
	QueryClassTrace(request QueryClassTraceRequest) (QueryClassTraceResponse, sdk.Error)
	QueryClassTraces(request QueryClassTracesRequest) (QueryClassTracesResponse, sdk.Error)
	QueryEscrowAddress(request QueryEscrowAddressRequest) (QueryEscrowAddressResponse, sdk.Error)
}