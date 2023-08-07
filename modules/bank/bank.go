package bank

import (
	"context"
	"fmt"

	sdk "github.com/irisnet/core-sdk-go/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type bankClient struct {
	sdk.BaseClient
	codec.Codec
	queryCli banktypes.QueryClient
}

// bank NewClient
func NewClient(bc sdk.BaseClient, cdc codec.Codec) Client {

	return bankClient{
		queryCli: banktypes.NewQueryClient(bc.GrpcConn()),

		Codec: cdc,
	}
}

func (b bankClient) Name() string {
	return banktypes.ModuleName
}

func (b bankClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	banktypes.RegisterInterfaces(registry)
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := b.BaseClient.QueryAccount(address)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	return account, nil
}

// TotalSupply queries the total supply of all coins.
func (b bankClient) TotalSupply() (types.Coins, sdk.Error) {

	resp, err := b.queryCli.TotalSupply(
		context.Background(),
		&banktypes.QueryTotalSupplyRequest{},
	)

	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return resp.Supply, nil
}

// Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount types.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s not found", baseTx.From)
	}

	amt, err := b.ToMinCoin(amount...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	outAddr, err1 := types.AccAddressFromBech32(to)
	if err1 != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", to))
	}

	msg := banktypes.NewMsgSend(sender, outAddr, amt)
	return b.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (b bankClient) SendWitchSpecAccountInfo(to string, sequence, accountNumber uint64, amount types.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s not found", baseTx.From)
	}

	amt, err := b.ToMinCoin(amount...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	outAddr, err1 := types.AccAddressFromBech32(to)
	if err1 != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", to))
	}

	msg := banktypes.NewMsgSend(sender, outAddr, amt)
	return b.BuildAndSendWithAccount(sender.String(), accountNumber, sequence, []types.Msg{msg}, baseTx)
}

func (b bankClient) MultiSend(request MultiSendRequest, baseTx sdk.BaseTx) (resTxs []sdk.ResultTx, err sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrapf("%s not found", baseTx.From)
	}

	var inputs = make([]banktypes.Input, len(request.Receipts))
	var outputs = make([]banktypes.Output, len(request.Receipts))
	for i, receipt := range request.Receipts {
		amt, err := b.ToMinCoin(receipt.Amount...)
		if err != nil {
			return nil, sdk.Wrap(err)
		}

		outAddr, e := types.AccAddressFromBech32(receipt.Address)
		if e != nil {
			return nil, sdk.Wrapf(fmt.Sprintf("%s invalid address", receipt.Address))
		}

		inputs[i] = banktypes.NewInput(sender, amt)
		outputs[i] = banktypes.NewOutput(outAddr, amt)
	}

	msg := banktypes.NewMsgMultiSend(inputs, outputs)
	res, err := b.BuildAndSend([]types.Msg{msg}, baseTx)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resTxs = append(resTxs, res)
	return
}
