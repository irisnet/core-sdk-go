package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	kmg "github.com/irisnet/core-sdk-go/crypto/keyring"
	"github.com/irisnet/core-sdk-go/crypto/keys/ethsecp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"
	"github.com/irisnet/core-sdk-go/store"
)

var _ IKeyManager = KeyManager{}

type IKeyManager interface {
	Sign(name, password string, data []byte) ([]byte, cryptotypes.PubKey, error)
	Insert(name, password string) (string, string, error)
	Recover(name, password, mnemonic, hdPath string) (string, error)
	Import(name, password string, privKeyArmor string) (address string, err error)
	Export(name, password string) (privKeyArmor string, err error)
	Delete(name, password string) error
	Find(name, password string) (cryptotypes.PubKey, types.AccAddress, error)
	Add(name, password string) (address string, mnemonic string, err Error)
}

type KeyManager struct {
	KeyDAO store.KeyDAO
	Algo   string
}

func NewKeyManager(keyDAO store.KeyDAO, algo string) KeyManager {
	return KeyManager{
		KeyDAO: keyDAO,
		Algo:   algo,
	}
}

func (k KeyManager) Add(name, password string) (string, string, Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, Wrap(err)
}

func (k KeyManager) Sign(name, password string, data []byte) ([]byte, cryptotypes.PubKey, error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	km, err := kmg.NewPrivateKeyManager([]byte(info.PrivKeyArmor), string(info.Algo))
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	signByte, err := km.Sign(data)
	if err != nil {
		return nil, nil, err
	}

	return signByte, km.ExportPubKey(), nil
}

func (k KeyManager) Insert(name, password string) (string, string, error) {
	if k.KeyDAO.Has(name) {
		return "", "", fmt.Errorf("name %s has existed", name)
	}

	km, err := kmg.NewAlgoKeyManager(k.Algo)
	if err != nil {
		return "", "", err
	}

	mnemonic, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       MarshalPubkey(pubKey),
		PrivKeyArmor: string(MarshalPrivKey(priv)),
		Algo:         k.Algo,
	}

	if err = k.KeyDAO.Write(name, password, info); err != nil {
		return "", "", err
	}
	return address, mnemonic, nil
}

func (k KeyManager) Recover(name, password, mnemonic, hdPath string) (string, error) {
	if k.KeyDAO.Has(name) {
		return "", fmt.Errorf("name %s has existed", name)
	}
	var (
		km  kmg.KeyManager
		err error
	)
	if hdPath == "" {
		km, err = kmg.NewMnemonicKeyManager(mnemonic, k.Algo)
	} else {
		km, err = kmg.NewMnemonicKeyManagerWithHDPath(mnemonic, k.Algo, hdPath)
	}

	if err != nil {
		return "", err
	}

	_, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       MarshalPubkey(pubKey),
		PrivKeyArmor: string(MarshalPrivKey(priv)),
		Algo:         k.Algo,
	}

	if err = k.KeyDAO.Write(name, password, info); err != nil {
		return "", err
	}

	return address, nil
}

func (k KeyManager) Import(name, password, armor string) (string, error) {
	if k.KeyDAO.Has(name) {
		return "", fmt.Errorf("%s has existed", name)
	}

	km := kmg.NewKeyManager()

	priv, _, err := km.ImportPrivKey(armor, password)
	if err != nil {
		return "", err
	}

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       MarshalPubkey(pubKey),
		PrivKeyArmor: string(MarshalPrivKey(priv)),
		Algo:         k.Algo,
	}

	err = k.KeyDAO.Write(name, password, info)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (k KeyManager) Export(name, password string) (armor string, err error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return armor, fmt.Errorf("name %s not exist", name)
	}

	km, err := kmg.NewPrivateKeyManager([]byte(info.PrivKeyArmor), info.Algo)
	if err != nil {
		return "", err
	}

	return km.ExportPrivKey(password)
}

func (k KeyManager) Delete(name, password string) error {
	return k.KeyDAO.Delete(name, password)
}

func (k KeyManager) Find(name, password string) (cryptotypes.PubKey, types.AccAddress, error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return nil, nil, WrapWithMessage(err, "name %s not exist", name)
	}

	pubKey, err := legacy.PubKeyFromBytes(info.PubKey)
	if err != nil {
		return nil, nil, WrapWithMessage(err, "name %s not exist", name)
	}
	return FromTmPubKey(info.Algo, pubKey), types.AccAddress(pubKey.Address().Bytes()), nil
}

func MarshalPubkey(pubKey cryptotypes.PubKey) []byte {
	return legacy.Cdc.MustMarshal(pubKey)
}

func MarshalPrivKey(privKey cryptotypes.PrivKey) []byte {
	return legacy.Cdc.MustMarshal(privKey)
}

func FromTmPubKey(Algo string, pubKey cryptotypes.PubKey) cryptotypes.PubKey {
	var pubkey cryptotypes.PubKey
	pubkeyBytes := pubKey.Bytes()
	switch Algo {
	case "sm2":
		pubkey = &sm2.PubKey{Key: pubkeyBytes}
	case "secp256k1":
		pubkey = &secp256k1.PubKey{Key: pubkeyBytes}
	case ethsecp256k1.KeyType:
		pubkey = &ethsecp256k1.PubKey{Key: pubkeyBytes}
	}
	return pubkey
}
