package feegrant

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/irisnet/core-sdk-go/common/codec/legacy"
	"github.com/irisnet/core-sdk-go/common/codec/types"

	sdk "github.com/irisnet/core-sdk-go/types"
)

const (
	ModuleName = "fee grant"

	TypeMsgGrantAllowance  = "grant allowance"
	TypeMsgRevokeAllowance = "revoke allowance"
)

var _, _ sdk.Msg = &MsgGrantAllowance{}, &MsgRevokeAllowance{}

// NewMsgGrantAllowance creates a new MsgGrantAllowance.
//nolint:interfacer
func NewMsgGrantAllowance(feeAllowance FeeAllowanceI, granter, grantee sdk.AccAddress) (*MsgGrantAllowance, error) {
	msg, ok := feeAllowance.(proto.Message)
	if !ok {
		return nil, errors.New("cannot proto marshal")
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &MsgGrantAllowance{
		Granter:   granter.String(),
		Grantee:   grantee.String(),
		Allowance: any,
	}, nil
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgGrantAllowance) ValidateBasic() error {
	if msg.Granter == "" {
		return errors.New("missing granter address")
	}
	if msg.Grantee == "" {
		return errors.New("missing grantee address")
	}
	if msg.Grantee == msg.Granter {
		return errors.New("cannot self-grant fee authorization")
	}

	allowance, err := msg.GetFeeAllowanceI()
	if err != nil {
		return err
	}

	return allowance.ValidateBasic()
}

// GetSigners gets the granter account associated with an allowance
func (msg MsgGrantAllowance) GetSigners() []sdk.AccAddress {
	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{granter}
}

// Type implements the LegacyMsg.Type method.
func (msg MsgGrantAllowance) Type() string {
	return MsgTypeURL(&msg)
}

// Route implements the LegacyMsg.Route method.
func (msg MsgGrantAllowance) Route() string {
	return MsgTypeURL(&msg)
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (msg MsgGrantAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

// GetFeeAllowanceI returns unpacked FeeAllowance
func (msg MsgGrantAllowance) GetFeeAllowanceI() (FeeAllowanceI, error) {
	allowance, ok := msg.Allowance.GetCachedValue().(FeeAllowanceI)
	if !ok {
		return nil, errors.New("failed to get allowance")
	}

	return allowance, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgGrantAllowance) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var allowance FeeAllowanceI
	return unpacker.UnpackAny(msg.Allowance, &allowance)
}

// NewMsgRevokeAllowance returns a message to revoke a fee allowance for a given
// granter and grantee
//nolint:interfacer
func NewMsgRevokeAllowance(granter sdk.AccAddress, grantee sdk.AccAddress) MsgRevokeAllowance {
	return MsgRevokeAllowance{Granter: granter.String(), Grantee: grantee.String()}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRevokeAllowance) ValidateBasic() error {
	if msg.Granter == "" {
		return errors.New("missing granter address")
	}
	if msg.Grantee == "" {
		return errors.New("missing grantee address")
	}
	if msg.Grantee == msg.Granter {
		return errors.New("addresses must be different")
	}

	return nil
}

// GetSigners gets the granter address associated with an Allowance
// to revoke.
func (msg MsgRevokeAllowance) GetSigners() []sdk.AccAddress {
	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{granter}
}

// Type implements the LegacyMsg.Type method.
func (msg MsgRevokeAllowance) Type() string {
	return MsgTypeURL(&msg)
}

// Route implements the LegacyMsg.Route method.
func (msg MsgRevokeAllowance) Route() string {
	return MsgTypeURL(&msg)
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (msg MsgRevokeAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

// MsgTypeURL returns the TypeURL of a `sdk.Msg`.
func MsgTypeURL(msg sdk.Msg) string {
	return "/" + proto.MessageName(msg)
}
