package perm

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	sdk "github.com/irisnet/core-sdk-go/types"
)

const (
	// ModuleName is the name of the perm module
	ModuleName         = "perm"
	TypeMsgAssignRoles = "assign_roles" // type for MsgAssignRoles
)

var _ sdk.Msg = &MsgAssignRoles{}

func NewMsgAssignRoles(operator, address sdk.AccAddress, roles []Role) (*MsgAssignRoles, error) {

	return &MsgAssignRoles{
		Operator: operator.String(),
		Address:  address.String(),
		Roles:    roles,
	}, nil
}

func (m MsgAssignRoles) Route() string {
	return ModuleName
}

func (m MsgAssignRoles) Type() string {
	return MsgTypeURL(&m)
}

func (m MsgAssignRoles) ValidateBasic() error {
	if len(m.Address) == 0 {
		return errors.New("address missing")
	}
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	if len(m.Roles) == 0 {
		return errors.New("roles missing")
	}
	return nil
}

func (m MsgAssignRoles) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgAssignRoles) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

// MsgTypeURL returns the TypeURL of a `sdk.Msg`.
func MsgTypeURL(msg sdk.Msg) string {
	return "/" + proto.MessageName(msg)
}
