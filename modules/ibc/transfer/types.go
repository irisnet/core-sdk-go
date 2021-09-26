package transfer

import (
	fmt "fmt"

	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
)

const (
	// ModuleName defines the IBC transfer name
	ModuleName = "transfer"
	// DenomPrefix is the prefix used for internal SDK coin representation.
	DenomPrefix = "ibc"
)

var _ types.Msg = &MsgTransfer{}

func (msg MsgTransfer) ValidateBasic() error {
	if _, err := types.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.Wrapf(ErrTodo, "invalid creator")
	}
	return nil
}

func (msg MsgTransfer) GetSigners() []types.AccAddress {
	creator, err := types.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []types.AccAddress{creator}
}

// String returns a string representation of Height
func (h Height) String() string {
	return fmt.Sprintf("%d-%d", h.RevisionNumber, h.RevisionHeight)
}
