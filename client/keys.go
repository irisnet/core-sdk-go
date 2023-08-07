package client

import (
	sdk "github.com/irisnet/core-sdk-go/types"
)

func (base *baseClient) Add(name, password string) (string, string, sdk.Error) {
	address, mnemonic, err := base.KeyManager.Insert(name, password)
	return address, mnemonic, sdk.Wrap(err)
}

func (base *baseClient) Recover(name, password, mnemonic string) (string, sdk.Error) {
	address, err := base.KeyManager.Recover(name, password, mnemonic, "")
	return address, sdk.Wrap(err)
}

func (base *baseClient) RecoverWithHDPath(name, password, mnemonic, hdPath string) (string, sdk.Error) {
	address, err := base.KeyManager.Recover(name, password, mnemonic, hdPath)
	return address, sdk.Wrap(err)
}

func (base *baseClient) Import(name, password, privKeyArmor string) (string, sdk.Error) {
	address, err := base.KeyManager.Import(name, password, privKeyArmor)
	return address, sdk.Wrap(err)
}

func (base *baseClient) Export(name, password string) (string, sdk.Error) {
	keystore, err := base.KeyManager.Export(name, password)
	return keystore, sdk.Wrap(err)
}

func (base *baseClient) Delete(name, password string) sdk.Error {
	err := base.KeyManager.Delete(name, password)
	return sdk.Wrap(err)
}

func (base *baseClient) Show(name, password string) (string, sdk.Error) {
	_, address, err := base.KeyManager.Find(name, password)
	if err != nil {
		return "", sdk.Wrap(err)
	}
	return address.String(), nil
}
