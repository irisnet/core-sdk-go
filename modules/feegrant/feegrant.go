package feegrant

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"

	sdk "github.com/irisnet/core-sdk-go/types"
)

type feeGrantClient struct {
	sdk.BaseClient
	codec.Codec
}

// NewClient grant NewClient
func NewClient(bc sdk.BaseClient, cdc codec.Codec) Client {
	return feeGrantClient{
		BaseClient: bc,
		Codec:      cdc,
	}
}

func (b feeGrantClient) Name() string {
	return feegranttypes.ModuleName
}

func (b feeGrantClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	feegranttypes.RegisterInterfaces(registry)
}

func (f feeGrantClient) GrantAllowance(
	granter, grantee types.AccAddress, feeAllowance feegranttypes.FeeAllowanceI, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	granter, err := types.AccAddressFromBech32(granter.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", granter))
	}

	grantee, err = types.AccAddressFromBech32(grantee.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", grantee))
	}

	msg, error := feegranttypes.NewMsgGrantAllowance(feeAllowance, granter, grantee)
	if error != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s", error))
	}
	return f.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (f feeGrantClient) RevokeAllowance(granter, grantee types.AccAddress, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	granter, err := types.AccAddressFromBech32(granter.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", granter))
	}

	grantee, err = types.AccAddressFromBech32(grantee.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", grantee))
	}

	msg := feegranttypes.NewMsgRevokeAllowance(granter, grantee)
	res, err := f.BuildAndSend([]types.Msg{&msg}, baseTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	return res, sdk.Wrap(err)
}
