package hd

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"
	sm22 "github.com/irisnet/core-sdk-go/crypto/keys/sm2"
)

const (
	// Sm2Type represents the Sm2Type signature system.
	Sm2Type = hd.PubKeyType("sm2")
)

var (
	_ keyring.SignatureAlgo = Sm2

	Sm2 = sm2Algo{}
)

type sm2Algo struct{}

func (s sm2Algo) Name() hd.PubKeyType {
	return Sm2Type
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s sm2Algo) Derive() hd.DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		masterPriv, ch := hd.ComputeMastersFromSeed(seed)
		if len(hdPath) == 0 {
			return masterPriv[:], nil
		}
		derivedKey, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath)
		return derivedKey[:], err
	}
}

// Generate generates a sm2 private key from the given bytes.
func (s sm2Algo) Generate() hd.GenerateFn {
	return func(bz []byte) cryptotypes.PrivKey {
		var bzArr [sm22.PrivKeySize]byte
		copy(bzArr[:], bz)
		return &sm22.PrivKey{Key: bzArr[:]}
	}
}
