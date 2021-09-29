package auth

import (
	"errors"

	"github.com/gogo/protobuf/proto"

	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/core-sdk-go/codec"
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	cryptotypes "github.com/irisnet/core-sdk-go/crypto/types"
	"github.com/irisnet/core-sdk-go/types"
)

//BaseAccount Have they all been implemented
var _ AccountI = (*BaseAccount)(nil)

// AccountI is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements AccountI.
type AccountI interface {
	proto.Message

	GetAddress() types.AccAddress
	SetAddress(types.AccAddress) error // errors if already set.

	GetPubKey() cryptotypes.PubKey // can return nil.
	SetPubKey(cryptotypes.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	// Ensure that account implements stringer
	String() string
}

// GetAddress - Implements types.AccountI.
func (acc BaseAccount) GetAddress() types.AccAddress {
	addr, _ := types.AccAddressFromBech32(acc.Address)
	return addr
}

// SetAddress - Implements types.AccountI.
func (acc *BaseAccount) SetAddress(addr types.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}

	acc.Address = addr.String()
	return nil
}

// GetPubKey - Implements types.AccountI.
func (acc BaseAccount) GetPubKey() (pk cryptotypes.PubKey) {
	if acc.PubKey == nil {
		return nil
	}
	content, ok := acc.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil
	}
	return content
}

// SetPubKey - Implements types.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey cryptotypes.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
		return nil
	}
	any, err := codectypes.NewAnyWithValue(pubKey)
	if err == nil {
		acc.PubKey = any
	}
	return err
}

// GetAccountNumber - Implements AccountI
func (acc BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements AccountI
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements types.AccountI.
func (acc BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements types.AccountI.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &acc)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

// Convert return a types.BaseAccount
// in order to unpack pubKey so not use Convert()
func (acc *BaseAccount) ConvertAccount(cdc codec.Codec) interface{} {
	account := types.BaseAccount{
		Address:       acc.Address,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	var pkStr string
	if acc.PubKey == nil {
		return account
	}

	var pk crypto.PubKey
	if err := cdc.UnpackAny(acc.PubKey, &pk); err != nil {
		return types.BaseAccount{}
	}

	pkStr, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pk)
	if err != nil {
		panic(err)
	}

	account.PubKey = pkStr
	return account
}
