package client

import (
	"context"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	cache "github.com/irisnet/core-sdk-go/cache"
	codec "github.com/irisnet/core-sdk-go/codec"
	"github.com/irisnet/core-sdk-go/modules/auth"
	"github.com/irisnet/core-sdk-go/modules/bank"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
)

// Must be used with locker, otherwise there are thread safety issues
type AccountQuery struct {
	types.Queries
	types.GRPCClient
	log.Logger
	cache.Cache
	cdc        codec.Codec
	Km         types.KeyManager
	expiration time.Duration
}

func (a AccountQuery) QueryAndRefreshAccount(address string) (types.BaseAccount, error) {
	account, err := a.Get(a.prefixKey(address))
	if err != nil {
		return a.refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := types.BaseAccount{
		Address:       address,
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	a.saveAccount(baseAcc)

	a.Debug("query account from cache", "address", address)
	return baseAcc, nil
}

func (a AccountQuery) QueryAccount(address string) (types.BaseAccount, error) {
	conn, err := a.GenConn()
	if err != nil {
		return types.BaseAccount{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	request := &auth.QueryAccountRequest{
		Address: address,
	}

	response, err := auth.NewQueryClient(conn).Account(context.Background(), request)
	if err != nil {
		return types.BaseAccount{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	var baseAccount auth.AccountI
	if err := a.cdc.UnpackAny(response.Account, &baseAccount); err != nil {
		return types.BaseAccount{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	account := baseAccount.(*auth.BaseAccount).ConvertAccount(a.cdc).(types.BaseAccount)

	breq := &bank.QueryAllBalancesRequest{
		Address:    address,
		Pagination: nil,
	}
	balances, err := bank.NewQueryClient(conn).AllBalances(context.Background(), breq)
	if err != nil {
		return types.BaseAccount{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	account.Coins = balances.Balances
	return account, nil
}

func (a AccountQuery) QueryAddress(name, password string) (types.AccAddress, error) {
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
		return address, errors.Wrap(errors.ErrTodo, err.Error())
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

func (a AccountQuery) refresh(address string) (types.BaseAccount, error) {
	account, err := a.QueryAccount(address)
	if err != nil {
		a.Error("update cache failed", "address", address, "errMsg", err.Error())
		return types.BaseAccount{}, errors.Wrap(errors.ErrTodo, err.Error())
	}

	a.saveAccount(account)
	return account, nil
}

func (a AccountQuery) saveAccount(account types.BaseAccount) {
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
