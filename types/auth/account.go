package auth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	commoncodec "github.com/irisnet/core-sdk-go/common/codec"
)

var (
	_ Account = (*EthAccount)(nil)
	//_ BaseAccount                 = (*EthAccount)(nil)
)

// ----------------------------------------------------------------------------
// Main Ethermint account
// ----------------------------------------------------------------------------

// ProtoAccount defines the prototype function for BaseAccount used for an
// AccountKeeper.
func ProtoAccount() Account {
	return &EthAccount{
		BaseAccount: &BaseAccount{},
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() common.Address {
	return common.BytesToAddress(acc.GetAddress().Bytes())
}

// GetCodeHash returns the account code hash in byte format
func (acc EthAccount) GetCodeHash() common.Hash {
	return common.HexToHash(acc.CodeHash)
}

type BaseAccountI interface {
	ConvertAccount(cdc commoncodec.Marshaler) interface{}
}
