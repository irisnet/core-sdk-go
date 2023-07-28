package bank

import (
	"github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

// expose bank module api for user
type Client interface {
	sdk.Module
	Send(to string, amount types.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	SendWitchSpecAccountInfo(to string, sequence, accountNumber uint64, amount types.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	MultiSend(receipts MultiSendRequest, baseTx sdk.BaseTx) ([]sdk.ResultTx, sdk.Error)
	QueryAccount(address string) (sdk.BaseAccount, sdk.Error)
	TotalSupply() (sdktypes.Coins, sdk.Error)
}

type Receipt struct {
	Address string         `json:"address"`
	Amount  types.DecCoins `json:"amount"`
}

type MultiSendRequest struct {
	Receipts []Receipt
}

func (msr MultiSendRequest) Len() int {
	return len(msr.Receipts)
}

func (msr MultiSendRequest) Sub(begin, end int) sdk.SplitAble {
	return MultiSendRequest{Receipts: msr.Receipts[begin:end]}
}
