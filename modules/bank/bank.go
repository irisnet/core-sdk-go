package bank

import (
	"context"
	"strings"

	"github.com/irisnet/core-sdk-go/codec"
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
)

type bankClient struct {
	types.BaseClient
	codec.Codec
}

// bank NewClient
func NewClient(bc types.BaseClient, cdc codec.Codec) Client {
	return bankClient{
		BaseClient: bc,
		Codec:      cdc,
	}
}

func (b bankClient) Name() string {
	return ModuleName
}

func (b bankClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (types.BaseAccount, error) {
	account, err := b.BaseClient.QueryAccount(address)
	if err != nil {
		return types.BaseAccount{}, errors.Wrap(ErrQueryAccount, err.Error())
	}

	return account, nil
}

//  TotalSupply queries the total supply of all coins.
func (b bankClient) TotalSupply() (types.Coins, error) {
	conn, err := b.GenConn()
	if err != nil {
		return nil, errors.Wrap(ErrGenConn, err.Error())
	}

	resp, err := NewQueryClient(conn).TotalSupply(
		context.Background(),
		&QueryTotalSupplyRequest{},
	)
	if err != nil {
		return nil, errors.Wrap(ErrQueryTotalSupply, err.Error())
	}
	return resp.Supply, nil
}

// Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount types.DecCoins, baseTx types.BaseTx) (types.ResultTx, error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrapf(err, "%s not found", baseTx.From)
	}

	amt, err := b.ToMinCoin(amount...)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrToMinCoin, err.Error())
	}

	outAddr, err := types.AccAddressFromBech32(to)
	if err != nil {
		return types.ResultTx{}, errors.Wrapf(err, "%s invalid address", to)
	}

	msg := NewMsgSend(sender, outAddr, amt)
	return b.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (b bankClient) SendWitchSpecAccountInfo(to string, sequence, accountNumber uint64, amount types.DecCoins, baseTx types.BaseTx) (types.ResultTx, error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrapf(err, "%s not found", baseTx.From)
	}

	amt, err := b.ToMinCoin(amount...)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrToMinCoin, err.Error())
	}

	outAddr, err := types.AccAddressFromBech32(to)
	if err != nil {
		return types.ResultTx{}, errors.Wrapf(err, "%s invalid address", to)
	}

	msg := NewMsgSend(sender, outAddr, amt)
	return b.BuildAndSendWithAccount(sender.String(), accountNumber, sequence, []types.Msg{msg}, baseTx)
}

func (b bankClient) MultiSend(request MultiSendRequest, baseTx types.BaseTx) (resTxs []types.ResultTx, err error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, errors.Wrapf(err, "%s not found", baseTx.From)
	}

	if len(request.Receipts) > maxMsgLen {
		return b.SendBatch(sender, request, baseTx)
	}

	var inputs = make([]Input, len(request.Receipts))
	var outputs = make([]Output, len(request.Receipts))
	for i, receipt := range request.Receipts {
		amt, err := b.ToMinCoin(receipt.Amount...)
		if err != nil {
			return nil, errors.Wrap(ErrToMinCoin, err.Error())
		}

		outAddr, e := types.AccAddressFromBech32(receipt.Address)
		if e != nil {
			return nil, errors.Wrapf(err, "%s invalid address", receipt.Address)
		}

		inputs[i] = NewInput(sender, amt)
		outputs[i] = NewOutput(outAddr, amt)
	}

	msg := NewMsgMultiSend(inputs, outputs)
	res, err := b.BuildAndSend([]types.Msg{msg}, baseTx)
	if err != nil {
		return nil, errors.Wrap(ErrBuildAndSend, err.Error())
	}

	resTxs = append(resTxs, res)
	return
}

func (b bankClient) SendBatch(sender types.AccAddress, request MultiSendRequest, baseTx types.BaseTx) ([]types.ResultTx, error) {
	batchReceipts := types.SubArray(maxMsgLen, request)

	var msgs types.Msgs
	for _, receipts := range batchReceipts {

		req := receipts.(MultiSendRequest)
		var inputs = make([]Input, len(req.Receipts))
		var outputs = make([]Output, len(req.Receipts))
		for i, receipt := range req.Receipts {
			amt, err := b.ToMinCoin(receipt.Amount...)
			if err != nil {
				return nil, errors.Wrap(ErrToMinCoin, err.Error())
			}

			outAddr, e := types.AccAddressFromBech32(receipt.Address)
			if e != nil {
				return nil, errors.Wrapf(err, "%s invalid address", receipt.Address)
			}

			inputs[i] = NewInput(sender, amt)
			outputs[i] = NewOutput(outAddr, amt)
		}
		msgs = append(msgs, NewMsgMultiSend(inputs, outputs))
	}
	return b.BaseClient.SendBatch(msgs, baseTx)
}

// SubscribeSendTx Subscribe MsgSend event and return subscription
func (b bankClient) SubscribeSendTx(from, to string, callback EventMsgSendCallback) types.Subscription {
	var builder = types.NewEventQueryBuilder()

	from = strings.TrimSpace(from)
	if len(from) != 0 {
		builder.AddCondition(types.NewCond(types.EventTypeMessage,
			types.AttributeKeySender).EQ(types.EventValue(from)))
	}

	to = strings.TrimSpace(to)
	if len(to) != 0 {
		builder.AddCondition(types.Cond("transfer.recipient").EQ(types.EventValue(to)))
	}

	subscription, _ := b.SubscribeTx(builder, func(data types.EventDataTx) {
		for _, msg := range data.Tx.GetMsgs() {
			if value, ok := msg.(*MsgSend); ok {
				callback(EventDataMsgSend{
					Height: data.Height,
					Hash:   data.Hash,
					From:   value.FromAddress,
					To:     value.ToAddress,
					Amount: value.Amount,
				})
			}
		}
	})
	return subscription
}
