package perm

import (
	"fmt"
	commoncodec "github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

type permClient struct {
	sdk.BaseClient
	commoncodec.Marshaler
}

// NewClient perm NewClient
func NewClient(bc sdk.BaseClient, cdc commoncodec.Marshaler) Client {
	return permClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (p permClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (p permClient) AssignRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	granter, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", granter))
	}

	operator, err := sdk.AccAddressFromBech32(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s invalid address", granter))
	}

	msg, error := NewMsgAssignRoles(operator, granter, roles)
	if error != nil {
		return sdk.ResultTx{}, sdk.Wrapf(fmt.Sprintf("%s", error))
	}
	return p.BuildAndSend([]sdk.Msg{msg}, baseTx)
}
