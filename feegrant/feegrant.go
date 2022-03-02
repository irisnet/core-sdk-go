package feegrant

import (
	"fmt"
	commoncodec "github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

type feeGrantClient struct {
	sdk.BaseClient
	commoncodec.Marshaler
}

// NewClient grant NewClient
func NewClient(bc sdk.BaseClient, cdc commoncodec.Marshaler) Client {
	return feeGrantClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (b feeGrantClient) Name() string {
	return ModuleName
}

func (b feeGrantClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (f feeGrantClient) GrantAllowance(granter, grantee sdk.AccAddress,feeAllowance FeeAllowanceI, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	granter, err := sdk.AccAddressFromBech32(granter.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", granter))
	}

	grantee, err = sdk.AccAddressFromBech32(grantee.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", grantee))
	}

	msg, error := NewMsgGrantAllowance(feeAllowance,granter, grantee)
	if error != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s", error))
	}
	return f.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (f feeGrantClient) RevokeAllowance(granter, grantee sdk.AccAddress, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	granter, err := sdk.AccAddressFromBech32(granter.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", granter))
	}

	grantee, err = sdk.AccAddressFromBech32(grantee.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", grantee))
	}

	msg := NewMsgRevokeAllowance(granter, grantee)
	res, err := f.BuildAndSend([]sdk.Msg{&msg}, baseTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	return res, sdk.Wrap(err)
}