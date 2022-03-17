package hd

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	ethcrypto"github.com/ethereum/go-ethereum/crypto"
	ethsecp256k1 "github.com/irisnet/core-sdk-go/common/crypto/keys/eth_secp256k1"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/core-sdk-go/common/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/sm2"
)

type SignatureAlgo interface {
	Name() PubKeyType
	Derive() DeriveFn
	Generate() GenerateFn
}

func NewSigningAlgoFromString(str string) (SignatureAlgo, error) {
	switch str {
	case string(Secp256k1.Name()):
		return Secp256k1, nil
	case string(Sm2.Name()):
		return Sm2, nil
	case string(EthSecp256k1.Name()):
		return EthSecp256k1, nil
	default:
		return nil, fmt.Errorf("provided algorithm `%s` is not supported", str)
	}
}

// PubKeyType defines an algorithm to derive key-pairs which can be used for cryptographic signing.
type PubKeyType string

const (
	// MultiType implies that a pubkey is a multisignature
	MultiType = PubKeyType("multi")
	// Secp256k1Type uses the Bitcoin secp256k1 ECDSA parameters.
	Secp256k1Type = PubKeyType("secp256k1")
	// Ed25519Type represents the Ed25519Type signature system.
	// It is currently not supported for end-user keys (wallets/ledgers).
	Ed25519Type = PubKeyType("ed25519")
	// Sr25519Type represents the Sr25519Type signature system.
	Sr25519Type = PubKeyType("sr25519")

	// Sm2Type represents the Sm2Type signature system.
	Sm2Type = PubKeyType("sm2")

	EthSecp256k1Type = PubKeyType("eth_secp256k1")
)

var (
	// Secp256k1 uses the Bitcoin secp256k1 ECDSA parameters.
	Secp256k1 = secp256k1Algo{}

	Sm2 = sm2Algo{}

	EthSecp256k1 = ethSecp256k1Algo{}
)

type DeriveFn func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error)
type GenerateFn func(bz []byte) crypto.PrivKey

type WalletGenerator interface {
	Derive(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error)
	Generate(bz []byte) crypto.PrivKey
}

type secp256k1Algo struct {
}

func (s secp256k1Algo) Name() PubKeyType {
	return Secp256k1Type
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s secp256k1Algo) Derive() DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		masterPriv, ch := ComputeMastersFromSeed(seed)
		if len(hdPath) == 0 {
			return masterPriv[:], nil
		}
		derivedKey, err := DerivePrivateKeyForPath(masterPriv, ch, hdPath)

		return derivedKey, err
	}
}

// Generate generates a secp256k1 private key from the given bytes.
func (s secp256k1Algo) Generate() GenerateFn {
	return func(bz []byte) crypto.PrivKey {
		var bzArr = make([]byte, secp256k1.PrivKeySize)
		copy(bzArr, bz)

		return &secp256k1.PrivKey{Key: bzArr}
	}
}

type sm2Algo struct{}

func (s sm2Algo) Name() PubKeyType {
	return Sm2Type
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s sm2Algo) Derive() DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		masterPriv, ch := ComputeMastersFromSeed(seed)
		if len(hdPath) == 0 {
			return masterPriv[:], nil
		}
		derivedKey, err := DerivePrivateKeyForPath(masterPriv, ch, hdPath)
		return derivedKey[:], err
	}
}

// Generate generates a sm2 private key from the given bytes.
func (s sm2Algo) Generate() GenerateFn {
	return func(bz []byte) crypto.PrivKey {
		var bzArr [sm2.PrivKeySize]byte
		copy(bzArr[:], bz)
		return &sm2.PrivKey{Key: bzArr[:]}
	}
}

type ethSecp256k1Algo struct{}

// Name returns eth_secp256k1
func (s ethSecp256k1Algo) Name() PubKeyType {
	return EthSecp256k1Type
}

// Derive derives and returns the eth_secp256k1 private key for the given mnemonic and HD path.
func (s ethSecp256k1Algo) Derive() DeriveFn {
	return func(mnemonic, bip39Passphrase, path string) ([]byte, error) {
		hdpath, err := accounts.ParseDerivationPath(path)
		if err != nil {
			return nil, err
		}

		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		// create a BTC-utils hd-derivation key chain
		masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
		if err != nil {
			return nil, err
		}

		key := masterKey
		for _, n := range hdpath {
			key, err = key.Derive(n)
			if err != nil {
				return nil, err
			}
		}

		// btc-utils representation of a secp256k1 private key
		privateKey, err := key.ECPrivKey()
		if err != nil {
			return nil, err
		}

		// cast private key to a convertible form (single scalar field element of secp256k1)
		// and then load into ethcrypto private key format.
		// TODO: add links to godocs of the two methods or implementations of them, to compare equivalency
		privateKeyECDSA := privateKey.ToECDSA()
		derivedKey := ethcrypto.FromECDSA(privateKeyECDSA)

		return derivedKey, nil
	}
}

// Generate generates a eth_secp256k1 private key from the given bytes.
func (s ethSecp256k1Algo) Generate() GenerateFn {
	return func(bz []byte) crypto.PrivKey {
		bzArr := make([]byte, ethsecp256k1.PrivKeySize)
		copy(bzArr, bz)
		return &ethsecp256k1.PrivKey{Key: bzArr}
	}
}
