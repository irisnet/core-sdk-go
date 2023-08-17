package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	rpc "github.com/tendermint/tendermint/rpc/client"
	"google.golang.org/grpc"
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
	AccountQuery
	TmQuery
}

type AccountQuery interface {
	QueryAccount(address string) (BaseAccount, Error)
	QueryAddress(name, password string) (types.AccAddress, Error)
}

type TmQuery interface {
	QueryTx(hash string) (*types.TxResponse, error)
	QueryTxs(events []string, page, limit int, orderBy string) (*types.SearchTxsResult, error)
	BlockMetadata(height int64) (BlockDetailMetadata, error)
}

type TokenManager interface {
	QueryToken(denom string) (Token, error)
	SaveTokens(tokens ...Token)
	ToMinCoin(coin ...types.DecCoin) (types.Coins, Error)
	ToMainCoin(coin ...types.Coin) (types.DecCoins, Error)
}

type Logger interface {
	SetLogger(log.Logger)
}

type KeyClient interface {
	Add(name, password string) (address string, mnemonic string, err Error)
	Recover(name, password, mnemonic string) (address string, err Error)
	RecoverWithHDPath(name, password, mnemonic, hdPath string) (address string, err Error)
	Import(name, password, privKeyArmor string) (address string, err Error)
	Export(name, password string) (privKeyArmor string, err Error)
	Delete(name, password string) Error
	Show(name, password string) (string, Error)
}

type BaseClient interface {
	TokenManager
	TxManager
	Queries
	Logger
	KeyClient
	rpc.Client

	GrpcConn() *grpc.ClientConn
}
