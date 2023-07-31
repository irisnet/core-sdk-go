package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

type TxManager interface {
	TmQuery
	BuildAndSend(msg []types.Msg, baseTx BaseTx) (ResultTx, Error)
	BuildAndSign(msg []types.Msg, baseTx BaseTx) ([]byte, Error)
	BuildTxHash(msg []types.Msg, baseTx BaseTx) (string, Error)
	BuildAndSendWithAccount(addr string, accountNumber, sequence uint64, msg []types.Msg, baseTx BaseTx) (ResultTx, Error)
	BuildAndSignWithAccount(addr string, accountNumber, sequence uint64, msg []types.Msg, baseTx BaseTx) ([]byte, Error)
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
	QueryParams(module string, res Response) Error
}

type StoreQuery interface {
	QueryWithResponse(path string, data interface{}, result Response) error
	Query(path string, data interface{}) ([]byte, error)
	QueryStore(key HexBytes, storeName string, height int64, prove bool) (abci.ResponseQuery, error)
}

type AccountQuery interface {
	QueryAccount(address string) (BaseAccount, Error)
	QueryAddress(name, password string) (types.AccAddress, Error)
}

type CacheManager interface {
	RemoveCache(address string) bool
}

type TmQuery interface {
	QueryTx(hash string) (*types.TxResponse, error)
	QueryTxs(events []string, page, limit int, orderBy string) (*types.SearchTxsResult, error)
	QueryBlock(height int64) (BlockDetail, error)
	BlockMetadata(height int64) (BlockDetailMetadata, error)
}

type TokenManager interface {
	QueryToken(denom string) (Token, error)
	SaveTokens(tokens ...Token)
	ToMinCoin(coin ...types.DecCoin) (types.Coins, Error)
	ToMainCoin(coin ...types.Coin) (types.DecCoins, Error)
}

type Logger interface {
	Logger() log.Logger
	SetLogger(log.Logger)
}

type BaseClient interface {
	TokenManager
	TxManager
	Queries
	TmClient
	Logger
	GRPCClient
	KeyManager
	CacheManager
}
