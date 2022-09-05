package feegrant

import (
	"github.com/gogo/protobuf/proto"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/pkg/errors"
)

var (
	_ types.UnpackInterfacesMessage = &Grant{}
)

// NewGrant creates a new FeeAllowanceGrant.
//nolint:interfacer
func NewGrant(granter, grantee sdk.AccAddress, feeAllowance FeeAllowanceI) (Grant, error) {
	msg, ok := feeAllowance.(proto.Message)
	if !ok {
		return Grant{}, errors.New("cannot proto marshal")
	}

	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return Grant{}, err
	}

	return Grant{
		Granter:   granter.String(),
		Grantee:   grantee.String(),
		Allowance: any,
	}, nil
}

// ValidateBasic performs basic validation on
// FeeAllowanceGrant
func (a Grant) ValidateBasic() error {
	if a.Granter == "" {
		return errors.New("missing granter address")
	}
	if a.Grantee == "" {
		return errors.New("missing grantee address")
	}
	if a.Grantee == a.Granter {
		return errors.New("cannot self-grant fee authorization")
	}

	f, err := a.GetGrant()
	if err != nil {
		return err
	}

	return f.ValidateBasic()
}

// GetGrant unpacks allowance
func (a Grant) GetGrant() (FeeAllowanceI, error) {
	allowance, ok := a.Allowance.GetCachedValue().(FeeAllowanceI)
	if !ok {
		return nil, errors.New("failed to get allowance")
	}

	return allowance, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (a Grant) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var allowance FeeAllowanceI
	return unpacker.UnpackAny(a.Allowance, &allowance)
}
