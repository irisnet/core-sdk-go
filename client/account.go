package client

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/core-sdk-go/common/cache"
	sdk "github.com/irisnet/core-sdk-go/types"
)

// Must be used with locker, otherwise there are thread safety issues
type AccountQuery struct {
	sdk.Queries
	sdk.GRPCClient
	log.Logger
	cache.Cache
	cdc        codec.Codec
	Km         sdk.KeyManager
	expiration time.Duration
}

func (a AccountQuery) QueryAndRefreshAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := a.Get(a.prefixKey(address))
	if err != nil {
		return a.refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := sdk.BaseAccount{
		Address:       address,
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	a.saveAccount(baseAcc)

	a.Debug("query account from cache", "address", address)
	return baseAcc, nil
}

func (a AccountQuery) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	conn, err := a.GenConn()

	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	request := &authtypes.QueryAccountRequest{
		Address: address,
	}

	response, err := authtypes.NewQueryClient(conn).Account(context.Background(), request)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	var baseAccount authtypes.AccountI
	if err := a.cdc.UnpackAny(response.Account, &baseAccount); err != nil {
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

func (a AccountQuery) QueryAddress(name, password string) (types.AccAddress, sdk.Error) {
	addr, err := a.Get(a.prefixKey(name))
	if err == nil {
		address, err := types.AccAddressFromBech32(addr.(string))
		if err != nil {
			a.Debug("invalid address", "name", name)
			_ = a.Remove(a.prefixKey(name))
		} else {
			return address, nil
		}
	}

	_, address, err := a.Km.Find(name, password)
	if err != nil {
		a.Debug("can't find account", "name", name)
		return address, sdk.Wrap(err)
	}

	if err := a.SetWithExpire(a.prefixKey(name), address.String(), a.expiration); err != nil {
		a.Debug("cache user failed", "name", name)
	}
	a.Debug("query user from cache", "name", name, "address", address.String())
	return address, nil
}

func (a AccountQuery) removeCache(address string) bool {
	return a.Remove(a.prefixKey(address))
}

func (a AccountQuery) refresh(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := a.QueryAccount(address)
	if err != nil {
		a.Error("update cache failed", "address", address, "errMsg", err.Error())
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	a.saveAccount(account)
	return account, nil
}

func (a AccountQuery) saveAccount(account sdk.BaseAccount) {
	address := account.Address
	info := accountInfo{
		N: account.AccountNumber,
		S: account.Sequence,
	}
	if err := a.SetWithExpire(a.prefixKey(address), info, a.expiration); err != nil {
		a.Debug("cache user failed", "address", account.Address)
		return
	}
	a.Debug("cache account", "address", address, "expiration", a.expiration.String())
}

func (a AccountQuery) prefixKey(address string) string {
	return fmt.Sprintf("account:%s", address)
}

type accountInfo struct {
	N uint64 `json:"n"`
	S uint64 `json:"s"`
}
