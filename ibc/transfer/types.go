package transfer

import (
	sdk "github.com/irisnet/core-sdk-go/types"
)

const (
	// ModuleName defines the IBC transfer name
	ModuleName = "transfer"
	// DenomPrefix is the prefix used for internal SDK coin representation.
	DenomPrefix = "ibc"
)

var (
	_ sdk.Msg = &MsgTransfer{}
)

func (msg MsgTransfer) Route() string {
	return ModuleName
}

func (msg MsgTransfer) Type() string {
	return "create_pool"
}

func (msg MsgTransfer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdk.Wrapf("invalid creator")
	}

	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
