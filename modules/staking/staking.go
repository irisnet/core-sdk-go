package staking

import (
	"context"

	"github.com/irisnet/core-sdk-go/codec"
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	cryptotypes "github.com/irisnet/core-sdk-go/crypto/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
	"github.com/irisnet/core-sdk-go/types/query"
)

type stakingClient struct {
	types.BaseClient
	codec.Codec
}

func NewClient(baseClient types.BaseClient, codec codec.Codec) Client {
	return &stakingClient{
		BaseClient: baseClient,
		Codec:      codec,
	}
}

func (sc stakingClient) Name() string {
	return ModuleName
}

func (sc stakingClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (sc stakingClient) CreateValidator(request CreateValidatorRequest, baseTx types.BaseTx) (types.ResultTx, error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}
	valAddr, err := types.ValAddressFromBech32(delegatorAddr.String())
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	values, err := sc.ToMinCoin(request.Value)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	// pk, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, request.Pubkey)
	// if err != nil {
	// 	return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	// }
	// pkAny, err := codectypes.PackAny(pk)
	// if err != nil {
	// 	return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	// }

	var pk cryptotypes.PubKey
	if err := sc.Codec.UnmarshalInterfaceJSON([]byte(request.Pubkey), &pk); err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	pkAny, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	msg := &MsgCreateValidator{
		Description:      Description{Moniker: request.Moniker},
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: valAddr.String(),
		Pubkey:           pkAny,
		Value:            values[0],
		Commission: CommissionRates{
			Rate:          request.Rate,
			MaxRate:       request.MaxRate,
			MaxChangeRate: request.MaxChangeRate,
		},
		MinSelfDelegation: request.MinSelfDelegation,
	}

	return sc.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (sc stakingClient) EditValidator(request EditValidatorRequest, baseTx types.BaseTx) (types.ResultTx, error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}
	valAddr, err := types.ValAddressFromBech32(delegatorAddr.String())
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	msg := &MsgEditValidator{
		Description: Description{
			Moniker:         request.Moniker,
			Identity:        request.Identity,
			Website:         request.Website,
			SecurityContact: request.SecurityContact,
			Details:         request.Details,
		},
		ValidatorAddress:  valAddr.String(),
		CommissionRate:    &request.CommissionRate,
		MinSelfDelegation: &request.MinSelfDelegation,
	}
	return sc.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (sc stakingClient) Delegate(request DelegateRequest, baseTx types.BaseTx) (types.ResultTx, error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	coins, err := sc.ToMinCoin(request.Amount)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	msg := &MsgDelegate{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: request.ValidatorAddr,
		Amount:           coins[0],
	}
	return sc.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (sc stakingClient) Undelegate(request UndelegateRequest, baseTx types.BaseTx) (types.ResultTx, error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	coins, err := sc.ToMinCoin(request.Amount)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}
	msg := &MsgUndelegate{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: request.ValidatorAddr,
		Amount:           coins[0],
	}
	return sc.BuildAndSend([]types.Msg{msg}, baseTx)
}

func (sc stakingClient) BeginRedelegate(request BeginRedelegateRequest, baseTx types.BaseTx) (types.ResultTx, error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}

	coins, err := sc.ToMinCoin(request.Amount)
	if err != nil {
		return types.ResultTx{}, errors.Wrap(ErrTodo, err.Error())
	}
	msg := &MsgBeginRedelegate{
		DelegatorAddress:    delegatorAddr.String(),
		ValidatorSrcAddress: request.ValidatorSrcAddress,
		ValidatorDstAddress: request.ValidatorDstAddress,
		Amount:              coins[0],
	}
	return sc.BuildAndSend([]types.Msg{msg}, baseTx)
}

// QueryValidators when status is "" will return all status' validator
// about status, you can see BondStatus_value
func (sc stakingClient) QueryValidators(status string, page, size uint64) (QueryValidatorsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(page, size)
	res, err := NewQueryClient(conn).Validators(
		context.Background(),
		&QueryValidatorsRequest{
			Status: status,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryValidatorsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert(sc.Codec).(QueryValidatorsResp), nil
}

func (sc stakingClient) QueryValidator(validatorAddr string) (QueryValidatorResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).Validator(
		context.Background(),
		&QueryValidatorRequest{
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryValidatorResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Validator.Convert(sc.Codec).(QueryValidatorResp), nil
}

func (sc stakingClient) QueryValidatorDelegations(validatorAddr string, page, size uint64) (QueryValidatorDelegationsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(page, size)
	res, err := NewQueryClient(conn).ValidatorDelegations(
		context.Background(),
		&QueryValidatorDelegationsRequest{
			ValidatorAddr: validatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryValidatorDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert().(QueryValidatorDelegationsResp), nil
}

func (sc stakingClient) QueryValidatorUnbondingDelegations(validatorAddr string, page, size uint64) (QueryValidatorUnbondingDelegationsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorUnbondingDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(page, size)
	res, err := NewQueryClient(conn).ValidatorUnbondingDelegations(
		context.Background(),
		&QueryValidatorUnbondingDelegationsRequest{
			ValidatorAddr: validatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryValidatorUnbondingDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert().(QueryValidatorUnbondingDelegationsResp), nil
}

func (sc stakingClient) QueryDelegation(delegatorAddr string, validatorAddr string) (QueryDelegationResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegationResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).Delegation(
		context.Background(),
		&QueryDelegationRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryDelegationResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.DelegationResponse.Convert().(QueryDelegationResp), nil
}

func (sc stakingClient) QueryUnbondingDelegation(delegatorAddr string, validatorAddr string) (QueryUnbondingDelegationResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryUnbondingDelegationResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).UnbondingDelegation(
		context.Background(),
		&QueryUnbondingDelegationRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryUnbondingDelegationResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Unbond.Convert().(QueryUnbondingDelegationResp), nil
}

func (sc stakingClient) QueryDelegatorDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorDelegationsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegatorDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(page, size)
	res, err := NewQueryClient(conn).DelegatorDelegations(
		context.Background(),
		&QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryDelegatorDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert().(QueryDelegatorDelegationsResp), nil
}

func (sc stakingClient) QueryDelegatorUnbondingDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorUnbondingDelegationsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegatorUnbondingDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(page, size)
	res, err := NewQueryClient(conn).DelegatorUnbondingDelegations(
		context.Background(),
		&QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: delegatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryDelegatorUnbondingDelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert().(QueryDelegatorUnbondingDelegationsResp), nil
}

func (sc stakingClient) QueryRedelegations(request QueryRedelegationsReq) (QueryRedelegationsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryRedelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(request.Page, request.Size)
	res, err := NewQueryClient(conn).Redelegations(
		context.Background(),
		&QueryRedelegationsRequest{
			DelegatorAddr:    request.DelegatorAddr,
			SrcValidatorAddr: request.SrcValidatorAddr,
			DstValidatorAddr: request.DstValidatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryRedelegationsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert().(QueryRedelegationsResp), nil
}

func (sc stakingClient) QueryDelegatorValidators(delegatorAddr string, page, size uint64) (QueryDelegatorValidatorsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegatorValidatorsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	offset, limit := types.ParsePage(page, size)
	res, err := NewQueryClient(conn).DelegatorValidators(
		context.Background(),
		&QueryDelegatorValidatorsRequest{
			DelegatorAddr: delegatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryDelegatorValidatorsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert(sc.Codec).(QueryDelegatorValidatorsResp), nil
}

func (sc stakingClient) QueryDelegatorValidator(delegatorAddr string, validatorAddr string) (QueryValidatorResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).DelegatorValidator(
		context.Background(),
		&QueryDelegatorValidatorRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryValidatorResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Validator.Convert(sc.Codec).(QueryValidatorResp), nil
}

// QueryHistoricalInfo tendermint only save latest 100 block, previous block is aborted
func (sc stakingClient) QueryHistoricalInfo(height int64) (QueryHistoricalInfoResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryHistoricalInfoResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).HistoricalInfo(
		context.Background(),
		&QueryHistoricalInfoRequest{
			Height: height,
		},
	)
	if err != nil {
		return QueryHistoricalInfoResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert(sc.Codec).(QueryHistoricalInfoResp), nil
}

func (sc stakingClient) QueryPool() (QueryPoolResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryPoolResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).Pool(
		context.Background(),
		&QueryPoolRequest{},
	)
	if err != nil {
		return QueryPoolResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	return QueryPoolResp{
		NotBondedTokens: res.Pool.NotBondedTokens,
		BondedTokens:    res.Pool.BondedTokens,
	}, nil
}

func (sc stakingClient) QueryParams() (QueryParamsResp, error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryParamsResp{}, errors.Wrap(ErrTodo, err.Error())
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{},
	)
	if err != nil {
		return QueryParamsResp{}, errors.Wrap(ErrTodo, err.Error())
	}
	return res.Convert().(QueryParamsResp), nil
}
