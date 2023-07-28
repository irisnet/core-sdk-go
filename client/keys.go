package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/legacy"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/irisnet/core-sdk-go/crypto/keys/ethsecp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"

	kmg "github.com/irisnet/core-sdk-go/crypto/keyring"

	"github.com/irisnet/core-sdk-go/store"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/irisnet/core-sdk-go/types"
)

var _ sdk.KeyManager = KeyManager{}

type KeyManager struct {
	KeyDAO store.KeyDAO
	Algo   string
}

func (k KeyManager) Add(name, password string) (string, string, sdk.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, sdk.Wrap(err)
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
		return nil, nil, sdk.WrapWithMessage(err, "name %s not exist", name)
	}

	pubKey, err := legacy.PubKeyFromBytes(info.PubKey)
	if err != nil {
		return nil, nil, sdk.WrapWithMessage(err, "name %s not exist", name)
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

type Client interface {
	Add(name, password string) (address string, mnemonic string, err sdk.Error)
	Recover(name, password, mnemonic string) (address string, err sdk.Error)
	RecoverWithHDPath(name, password, mnemonic, hdPath string) (address string, err sdk.Error)
	Import(name, password, privKeyArmor string) (address string, err sdk.Error)
	Export(name, password string) (privKeyArmor string, err sdk.Error)
	Delete(name, password string) sdk.Error
	Show(name, password string) (string, sdk.Error)
}

type keysClient struct {
	hd.BIP44Params
	sdk.KeyManager
}

func NewKeysClient(cfg sdk.ClientConfig, keyManager sdk.KeyManager) Client {
	BIP44Params, err := hd.NewParamsFromPath(cfg.BIP44Path)
	if err != nil {
		panic(err)
	}
	return keysClient{*BIP44Params, keyManager}
}

func (k keysClient) Add(name, password string) (string, string, sdk.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, sdk.Wrap(err)
}

func (k keysClient) Recover(name, password, mnemonic string) (string, sdk.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic, "")
	return address, sdk.Wrap(err)
}

func (k keysClient) RecoverWithHDPath(name, password, mnemonic, hdPath string) (string, sdk.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic, hdPath)
	return address, sdk.Wrap(err)
}

func (k keysClient) Import(name, password, privKeyArmor string) (string, sdk.Error) {
	address, err := k.KeyManager.Import(name, password, privKeyArmor)
	return address, sdk.Wrap(err)
}

func (k keysClient) Export(name, password string) (string, sdk.Error) {
	keystore, err := k.KeyManager.Export(name, password)
	return keystore, sdk.Wrap(err)
}

func (k keysClient) Delete(name, password string) sdk.Error {
	err := k.KeyManager.Delete(name, password)
	return sdk.Wrap(err)
}

func (k keysClient) Show(name, password string) (string, sdk.Error) {
	_, address, err := k.KeyManager.Find(name, password)
	if err != nil {
		return "", sdk.Wrap(err)
	}
	return address.String(), nil
}
