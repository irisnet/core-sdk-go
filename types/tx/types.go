package tx

import (
	"errors"
	"fmt"

	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/types"
)

// MaxGasWanted defines the max gas allowed.
const MaxGasWanted = uint64((1 << 63) - 1)

var _ codectypes.UnpackInterfacesMessage = &TxBody{}
var _ codectypes.UnpackInterfacesMessage = &Tx{}
var _ types.Tx = &Tx{}

// GetMsgs implements the GetMsgs method on types.Tx.
func (t *Tx) GetMsgs() []types.Msg {
	if t == nil || t.Body == nil {
		return nil
	}

	anys := t.Body.Messages
	res := make([]types.Msg, len(anys))
	for i, any := range anys {
		msg := any.GetCachedValue().(types.Msg)
		res[i] = msg
	}
	return res
}

// ValidateBasic implements the ValidateBasic method on types.Tx.
func (t *Tx) ValidateBasic() error {
	if t == nil {
		return fmt.Errorf("bad Tx")
	}

	body := t.Body
	if body == nil {
		return fmt.Errorf("missing TxBody")
	}

	authInfo := t.AuthInfo
	if authInfo == nil {
		return fmt.Errorf("missing AuthInfo")
	}

	fee := authInfo.Fee
	if fee == nil {
		return fmt.Errorf("missing fee")
	}

	if fee.GasLimit > MaxGasWanted {
		return fmt.Errorf(
			"invalid gas supplied; %d > %d", fee.GasLimit, MaxGasWanted,
		)
	}

	if fee.Amount.IsAnyNegative() {
		return fmt.Errorf(
			"invalid fee provided: %s", fee.Amount,
		)
	}

	sigs := t.Signatures

	if len(sigs) == 0 {
		return errors.New("no signatures supplied")
	}

	if len(sigs) != len(t.GetSigners()) {
		return fmt.Errorf(
			"wrong number of signers; expected %d, got %d", t.GetSigners(), len(sigs),
		)
	}

	return nil
}

// GetSigners retrieves all the signers of a tx.
func (t *Tx) GetSigners() []types.AccAddress {
	var signers []types.AccAddress
	seen := map[string]bool{}

	for _, msg := range t.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	return signers
}

// UnpackInterfaces implements the UnpackInterfaceMessages.UnpackInterfaces method
func (t *Tx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if t.Body != nil {
		return t.Body.UnpackInterfaces(unpacker)
	}
	return nil
}

// UnpackInterfaces implements the UnpackInterfaceMessages.UnpackInterfaces method
func (m *TxBody) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, any := range m.Messages {
		var msg types.Msg
		if err := unpacker.UnpackAny(any, &msg); err != nil {
			return err
		}
	}
	return nil
}

// RegisterInterfaces registers the types.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.tx.v1beta1.Tx", (*types.Tx)(nil))
	registry.RegisterImplementations((*types.Tx)(nil), &Tx{})
}
