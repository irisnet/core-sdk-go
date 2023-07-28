package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
)

var (
	SupportedAlgorithms = keyring.SigningAlgoList{hd.Secp256k1}
)
