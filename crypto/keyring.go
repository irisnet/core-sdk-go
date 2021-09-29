package crypto

import (
	"fmt"
	"strings"

	"github.com/cosmos/go-bip39"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/core-sdk-go/codec/legacy"
	"github.com/irisnet/core-sdk-go/crypto/hd"
	"github.com/irisnet/core-sdk-go/types/errors"
)

const defaultBIP39Passphrase = ""

type Keyring interface {
	Generate() (string, tmcrypto.PrivKey)
	Sign(data []byte) ([]byte, error)

	ImportPrivKey(armor, passphrase string) (tmcrypto.PrivKey, string, error)
	ExportPrivKey(password string) (armor string, err error)

	ExportPubKey() tmcrypto.PubKey
}

type keyring struct {
	privKey  tmcrypto.PrivKey
	mnemonic string
	algo     string
}

func NewKeyring() Keyring {
	return &keyring{}
}

func NewAlgoKeyring(algo string) (Keyring, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return NewMnemonicKeyring(mnemonic, algo)
}

func NewMnemonicKeyring(mnemonic string, algo string) (Keyring, error) {
	k := keyring{
		mnemonic: mnemonic,
		algo:     algo,
	}
	err := k.recoveryFromMnemonic(mnemonic, hd.FullPath, algo)
	return &k, err
}

func NewMnemonicKeyringWithHDPath(mnemonic, algo, hdPath string) (Keyring, error) {
	k := keyring{
		mnemonic: mnemonic,
		algo:     algo,
	}
	err := k.recoveryFromMnemonic(mnemonic, hdPath, algo)
	return &k, err
}

func NewPrivateKeyring(priv []byte, algo string) (Keyring, error) {
	privKey, err := legacy.PrivKeyFromBytes(priv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt private key")
	}
	k := keyring{
		privKey: privKey,
		algo:    algo,
	}
	return &k, err
}

func (m *keyring) Generate() (string, tmcrypto.PrivKey) {
	return m.mnemonic, m.privKey
}

func (m *keyring) Sign(data []byte) ([]byte, error) {
	return m.privKey.Sign(data)
}

func (m *keyring) recoveryFromMnemonic(mnemonic, hdPath, algoStr string) error {
	words := strings.Split(mnemonic, " ")
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic length should either be 12 or 24")
	}

	algo, err := hd.NewSigningAlgoFromString(algoStr)
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

func (m *keyring) ExportPrivKey(password string) (armor string, err error) {
	return EncryptArmorPrivKey(m.privKey, password, m.algo), nil
}

func (m *keyring) ImportPrivKey(armor, passphrase string) (tmcrypto.PrivKey, string, error) {

	privKey, algo, err := UnarmorDecryptPrivKey(armor, passphrase)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to decrypt private key")
	}

	m.privKey = privKey
	m.algo = algo
	return privKey, algo, nil
}

func (m *keyring) ExportPubKey() tmcrypto.PubKey {
	return m.privKey.PubKey()
}
