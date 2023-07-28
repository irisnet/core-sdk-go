package keyring

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec/legacy"

	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/pkg/errors"

	"github.com/cosmos/go-bip39"
)

const (
	defaultBIP39Passphrase = ""
)

type KeyManager interface {
	Generate() (string, cryptotypes.PrivKey)
	Sign(data []byte) ([]byte, error)

	ExportPrivKey(password string) (armor string, err error)
	ImportPrivKey(armor, passphrase string) (cryptotypes.PrivKey, string, error)

	ExportPubKey() cryptotypes.PubKey
}

type keyManager struct {
	privKey        cryptotypes.PrivKey
	mnemonic, algo string
}

func NewKeyManager() KeyManager {
	return &keyManager{}
}

func NewAlgoKeyManager(algo string) (KeyManager, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return NewMnemonicKeyManager(mnemonic, algo)
}

func NewMnemonicKeyManager(mnemonic string, algo string) (KeyManager, error) {
	k := keyManager{
		mnemonic: mnemonic,
		algo:     algo,
	}
	err := k.recoveryFromMnemonic(mnemonic, FullPath, algo)
	return &k, err
}

func NewMnemonicKeyManagerWithHDPath(mnemonic, algo, hdPath string) (KeyManager, error) {
	k := keyManager{
		mnemonic: mnemonic,
		algo:     algo,
	}
	err := k.recoveryFromMnemonic(mnemonic, hdPath, algo)
	return &k, err
}

func NewPrivateKeyManager(priv []byte, algo string) (KeyManager, error) {

	privKey, err := legacy.PrivKeyFromBytes(priv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt private key")
	}
	k := keyManager{
		privKey: privKey,
		algo:    algo,
	}
	return &k, err
}

func (m *keyManager) Generate() (string, cryptotypes.PrivKey) {
	return m.mnemonic, m.privKey
}

func (m *keyManager) Sign(data []byte) ([]byte, error) {
	return m.privKey.Sign(data)
}

func (m *keyManager) recoveryFromMnemonic(mnemonic, hdPath, algoStr string) error {
	words := strings.Split(mnemonic, " ")
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic length should either be 12 or 24")
	}

	algo, err := keyring.NewSigningAlgoFromString(algoStr, SupportedAlgorithms)
	if err != nil {
		return err
	}

	// create master key and derive first key for keyring
	derivedPriv, err := algo.Derive()(mnemonic, defaultBIP39Passphrase, hdPath)
	if err != nil {
		return err
	}

	privKey := algo.Generate()(derivedPriv)
	m.privKey = privKey
	m.algo = algoStr
	return nil
}

func (m *keyManager) ExportPrivKey(password string) (armor string, err error) {
	return crypto.EncryptArmorPrivKey(m.privKey, password, m.algo), nil
}

func (m *keyManager) ImportPrivKey(armor, passphrase string) (cryptotypes.PrivKey, string, error) {
	privKey, algo, err := crypto.UnarmorDecryptPrivKey(armor, passphrase)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to decrypt private key")
	}

	m.privKey = privKey
	m.algo = algo
	return privKey, algo, nil
}

func (m *keyManager) ExportPubKey() cryptotypes.PubKey {
	return m.privKey.PubKey()
}
