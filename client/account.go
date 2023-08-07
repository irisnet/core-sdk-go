package client

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/irisnet/core-sdk-go/types"
)

func (base *baseClient) QueryAndRefreshAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := base.Get(base.prefixKey(address))
	if err != nil {
		return base.refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := sdk.BaseAccount{
		Address:       address,
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	base.saveAccount(baseAcc)

	base.Debug("query account from cache", "address", address)
	return baseAcc, nil
}

func (base *baseClient) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	request := &authtypes.QueryAccountRequest{
		Address: address,
	}

	response, err := authtypes.NewQueryClient(base.grpcConn).Account(context.Background(), request)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	var baseAccount authtypes.AccountI
	if err := base.encodingConfig.Marshaler.UnpackAny(response.Account, &baseAccount); err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	account := sdk.BaseAccount{
		Address:       baseAccount.GetAddress().String(),
		AccountNumber: baseAccount.GetAccountNumber(),
		Sequence:      baseAccount.GetSequence(),
	}

	if baseAccount.GetPubKey() != nil {
		account.PubKey = baseAccount.GetPubKey().String()
	}

	return account, nil
}

func (base *baseClient) QueryAddress(name, password string) (types.AccAddress, sdk.Error) {
	addr, err := base.Get(base.prefixKey(name))
	if err == nil {
		address, err := types.AccAddressFromBech32(addr.(string))
		if err != nil {
			base.Debug("invalid address", "name", name)
			_ = base.Remove(base.prefixKey(name))
		} else {
			return address, nil
		}
	}

	_, address, err := base.KeyManager.Find(name, password)
	if err != nil {
		base.Debug("can't find account", "name", name)
		return address, sdk.Wrap(err)
	}

	if err := base.SetWithExpire(base.prefixKey(name), address.String(), base.expiration); err != nil {
		base.Debug("cache user failed", "name", name)
	}
	base.Debug("query user from cache", "name", name, "address", address.String())
	return address, nil
}

func (base *baseClient) removeCache(address string) bool {
	return base.Remove(base.prefixKey(address))
}

func (base *baseClient) refresh(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := base.QueryAccount(address)
	if err != nil {
		base.Error("update cache failed", "address", address, "errMsg", err.Error())
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	base.saveAccount(account)
	return account, nil
}

func (base *baseClient) saveAccount(account sdk.BaseAccount) {
	address := account.Address
	info := accountInfo{
		N: account.AccountNumber,
		S: account.Sequence,
	}
	if err := base.SetWithExpire(base.prefixKey(address), info, base.expiration); err != nil {
		base.Debug("cache user failed", "address", account.Address)
		return
	}
	base.Debug("cache account", "address", address, "expiration", base.expiration.String())
}

func (base *baseClient) prefixKey(address string) string {
	return fmt.Sprintf("account:%s", address)
}

type accountInfo struct {
	N uint64 `json:"n"`
	S uint64 `json:"s"`
}
