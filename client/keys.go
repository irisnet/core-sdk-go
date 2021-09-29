package client

import (
	"fmt"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/core-sdk-go/codec/legacy"
	"github.com/irisnet/core-sdk-go/crypto"
	"github.com/irisnet/core-sdk-go/crypto/hd"
	"github.com/irisnet/core-sdk-go/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"
	cryptotypes "github.com/irisnet/core-sdk-go/crypto/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
	"github.com/irisnet/core-sdk-go/types/store"
)

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

func (k KeyManager) Add(name, password string) (string, string, error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, errors.Wrap(errors.ErrTodo, err.Error())
}

func (k KeyManager) Sign(name, password string, data []byte) ([]byte, tmcrypto.PubKey, error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	km, err := crypto.NewPrivateKeyring([]byte(info.PrivKeyArmor), string(info.Algo))
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	signByte, err := km.Sign(data)
	if err != nil {
		return nil, nil, err
	}

	return signByte, FromTmPubKey(info.Algo, km.ExportPubKey()), nil
}

func (k KeyManager) Insert(name, password string) (string, string, error) {
	if k.KeyDAO.Has(name) {
		return "", "", fmt.Errorf("name %s has existed", name)
	}

	km, err := crypto.NewAlgoKeyring(k.Algo)
	if err != nil {
		return "", "", err
	}

	mnemonic, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       legacy.MarshalPubkey(pubKey),
		PrivKeyArmor: string(legacy.MarshalPrivKey(priv)),
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
		km  crypto.Keyring
		err error
	)
	if hdPath == "" {
		km, err = crypto.NewMnemonicKeyring(mnemonic, k.Algo)
	} else {
		km, err = crypto.NewMnemonicKeyringWithHDPath(mnemonic, k.Algo, hdPath)
	}

	if err != nil {
		return "", err
	}

	_, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       legacy.MarshalPubkey(pubKey),
		PrivKeyArmor: string(legacy.MarshalPrivKey(priv)),
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

	km := crypto.NewKeyring()

	priv, _, err := km.ImportPrivKey(armor, password)
	if err != nil {
		return "", err
	}

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       legacy.MarshalPubkey(pubKey),
		PrivKeyArmor: string(legacy.MarshalPrivKey(priv)),
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

	km, err := crypto.NewPrivateKeyring([]byte(info.PrivKeyArmor), info.Algo)
	if err != nil {
		return "", err
	}

	return km.ExportPrivKey(password)
}

func (k KeyManager) Delete(name, password string) error {
	return k.KeyDAO.Delete(name, password)
}

func (k KeyManager) Find(name, password string) (tmcrypto.PubKey, types.AccAddress, error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "name %s not exist", name)
	}

	pubKey, err := legacy.PubKeyFromBytes(info.PubKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "name %s not exist", name)
	}

	return FromTmPubKey(info.Algo, pubKey), types.AccAddress(pubKey.Address().Bytes()), nil
}

func FromTmPubKey(Algo string, pubKey tmcrypto.PubKey) cryptotypes.PubKey {
	var pubkey cryptotypes.PubKey
	pubkeyBytes := pubKey.Bytes()
	switch Algo {
	case "sm2":
		pubkey = &sm2.PubKey{Key: pubkeyBytes}
	case "secp256k1":
		pubkey = &secp256k1.PubKey{Key: pubkeyBytes}
	}
	return pubkey
}

type Client interface {
	Add(name, password string) (address string, mnemonic string, err error)
	Recover(name, password, mnemonic string) (address string, err error)
	RecoverWithHDPath(name, password, mnemonic, hdPath string) (address string, err error)
	Import(name, password, privKeyArmor string) (address string, err error)
	Export(name, password string) (privKeyArmor string, err error)
	Delete(name, password string) error
	Show(name, password string) (string, error)
}

type keysClient struct {
	hd.BIP44Params
	types.KeyManager
}

func NewKeysClient(cfg types.ClientConfig, keyManager types.KeyManager) Client {
	BIP44Params, err := hd.NewParamsFromPath(cfg.BIP44Path)
	if err != nil {
		panic(err)
	}
	return keysClient{(*BIP44Params), keyManager}
}

func (k keysClient) Add(name, password string) (string, string, error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, errors.Wrap(errors.ErrTodo, err.Error())
}

func (k keysClient) Recover(name, password, mnemonic string) (string, error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic, "")
	return address, errors.Wrap(errors.ErrTodo, err.Error())
}

func (k keysClient) RecoverWithHDPath(name, password, mnemonic, hdPath string) (string, error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic, hdPath)
	return address, errors.Wrap(errors.ErrTodo, err.Error())
}

func (k keysClient) Import(name, password, privKeyArmor string) (string, error) {
	address, err := k.KeyManager.Import(name, password, privKeyArmor)
	return address, errors.Wrap(errors.ErrTodo, err.Error())
}

func (k keysClient) Export(name, password string) (string, error) {
	keystore, err := k.KeyManager.Export(name, password)
	return keystore, errors.Wrap(errors.ErrTodo, err.Error())
}

func (k keysClient) Delete(name, password string) error {
	err := k.KeyManager.Delete(name, password)
	return errors.Wrap(errors.ErrTodo, err.Error())
}

func (k keysClient) Show(name, password string) (string, error) {
	_, address, err := k.KeyManager.Find(name, password)
	if err != nil {
		return "", errors.Wrap(errors.ErrTodo, err.Error())
	}
	return address.String(), nil
}
