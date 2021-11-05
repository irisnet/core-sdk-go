package bank

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/types"
)

// expose bank module api for user
type Client interface {
	types.Module
	// Send defines a method for sending coins from one account to another account.
	Send(to string, amount types.DecCoins, baseTx types.BaseTx) (ctypes.ResultTx, error)
	// MultiSend defines a method for sending coins from some accounts to other accounts.
	SendWitchSpecAccountInfo(to string, sequence, accountNumber uint64, amount types.DecCoins, baseTx types.BaseTx) (ctypes.ResultTx, error)
	MultiSend(receipts MultiSendRequest, baseTx types.BaseTx) ([]ctypes.ResultTx, error)
	SubscribeSendTx(from, to string, callback EventMsgSendCallback) types.Subscription
	QueryAccount(address string) (types.BaseAccount, error)
	TotalSupply() (types.Coins, error)
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

func (msr MultiSendRequest) Sub(begin, end int) types.SplitAble {
	return MultiSendRequest{Receipts: msr.Receipts[begin:end]}
}

type EventDataMsgSend struct {
	Height int64        `json:"height"`
	Hash   string       `json:"hash"`
	From   string       `json:"from"`
	To     string       `json:"to"`
	Amount []types.Coin `json:"amount"`
}

type EventMsgSendCallback func(EventDataMsgSend)
