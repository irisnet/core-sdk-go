package hd

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

var (
	// SupportedAlgorithms defines the list of signing algorithms used on Ethermint:
	//  - eth_secp256k1 (Ethereum)
	//  - secp256k1 (Tendermint)
	// - sm2
	SupportedAlgorithms = keyring.SigningAlgoList{EthSecp256k1, hd.Secp256k1, Sm2}
	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Ethermint for the Ledger device:
	//  - eth_secp256k1 (Ethereum)
	//  - secp256k1 (Tendermint)
	// - sm2
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{EthSecp256k1, hd.Secp256k1, Sm2}
)
