package types

import (
	grpc1 "github.com/gogo/protobuf/grpc"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type TxManager interface {
	TmQuery
	SendBatch(msgs Msgs, baseTx BaseTx) ([]ctypes.ResultTx, error)
	BuildAndSend(msg []Msg, baseTx BaseTx) (ctypes.ResultTx, error)
	BuildAndSign(msg []Msg, baseTx BaseTx) ([]byte, error)
	BuildTxHash(msg []Msg, baseTx BaseTx) (string, error)
	BuildAndSendWithAccount(addr string, accountNumber, sequence uint64, msg []Msg, baseTx BaseTx) (ctypes.ResultTx, error)
}

type Queries interface {
	StoreQuery
	AccountQuery
	TmQuery
}

type GRPCClient interface {
	GenConn() (grpc1.ClientConn, error)
}

type ParamQuery interface {
	QueryParams(module string, res Response) error
}

type StoreQuery interface {
	QueryWithResponse(path string, data interface{}, result Response) error
	Query(path string, data interface{}) ([]byte, error)
	QueryStore(key HexBytes, storeName string, height int64, prove bool) (abci.ResponseQuery, error)
}

type AccountQuery interface {
	QueryAccount(address string) (BaseAccount, error)
	QueryAddress(name, password string) (AccAddress, error)
}

type TmQuery interface {
	QueryTx(hash string) (ctypes.ResultTx, error)
	QueryTxs(builder *EventQueryBuilder, page, size *int) (ctypes.ResultTxSearch, error)
	QueryBlock(height int64) (BlockDetail, error)
}

type TokenManager interface {
	QueryToken(denom string) (Token, error)
	SaveTokens(tokens ...Token)
	ToMinCoin(coin ...DecCoin) (Coins, error)
	ToMainCoin(coin ...Coin) (DecCoins, error)
}

type Logger interface {
	Logger() log.Logger
	SetLogger(log.Logger)
}

type WSClient interface {
	SubscribeNewBlock(builder *EventQueryBuilder, handler EventNewBlockHandler) (Subscription, error)
	SubscribeTx(builder *EventQueryBuilder, handler EventTxHandler) (Subscription, error)
	SubscribeNewBlockHeader(handler EventNewBlockHeaderHandler) (Subscription, error)
	SubscribeValidatorSetUpdates(handler EventValidatorSetUpdatesHandler) (Subscription, error)
	Unsubscribe(subscription Subscription) error
}

type TmClient interface {
	ABCIClient
	SignClient
	WSClient
	StatusClient
	NetworkClient
}

type BaseClient interface {
	TxConfig() TxConfig
	TokenManager
	TxManager
	Queries
	TmClient
	Logger
	GRPCClient
	KeyManager
}
