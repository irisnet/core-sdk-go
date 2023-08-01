package client

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/tendermint/tendermint/libs/log"
	rpc "github.com/tendermint/tendermint/rpc/client"
	rpcclienthttp "github.com/tendermint/tendermint/rpc/client/http"
)

type rpcClient struct {
	rpc.Client
	log.Logger
	cdc       *codec.LegacyAmino
	txDecoder types.TxDecoder
}

func NewRPCClient(cfg sdk.ClientConfig,
	cdc *codec.LegacyAmino,
	txDecoder types.TxDecoder,
	logger log.Logger,
) sdk.TmClient {
	if len(cfg.WSAddr) == 0 {
		cfg.WSAddr = "/websocket"
	}
	client, err := rpcclienthttp.NewWithTimeout(
		cfg.RPCAddr,
		cfg.WSAddr,
		cfg.Timeout)
	if err != nil {
		panic(err)
	}

	//if err := client.Start(); err != nil {
	//	panic(err)
	//}
	return rpcClient{
		Client:    client,
		Logger:    logger,
		cdc:       cdc,
		txDecoder: txDecoder,
	}
}
