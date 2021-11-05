package staking

import (
	"github.com/irisnet/core-sdk-go/types/errors"
)

const Codespace = ModuleName

var (
	ErrInvalidAmount                      = errors.Register(Codespace, 1, "invalid amount")
	ErrInvalidMinSelfDelegation           = errors.Register(Codespace, 2, "invalid MinSelfDelegation")
	ErrInvalidCommission                  = errors.Register(Codespace, 3, "invalid commission")
	ErrInvalidDescription                 = errors.Register(Codespace, 4, "invalid description")
	ErrQueryAddress                       = errors.Register(Codespace, 5, "query address error")
	ErrToMinCoin                          = errors.Register(Codespace, 6, "ToMinCoin error")
	ErrNewAnyWithValue                    = errors.Register(Codespace, 7, "NewAnyWithValue error")
	ErrQueryValidator                     = errors.Register(Codespace, 8, " query validator error")
	ErrQueryValidatorDelegations          = errors.Register(Codespace, 9, "query validator delegations error")
	ErrQueryValidatorUnbondingDelegations = errors.Register(Codespace, 10, " query validator unbonding delegations error")
	ErrQueryDelegation                    = errors.Register(Codespace, 12, "query delegation error")
	ErrQueryUnbondingDelegation           = errors.Register(Codespace, 13, "query unbonding delegation error")
	ErrQueryDelegatorDelegations          = errors.Register(Codespace, 14, "query delegator delegations error")
	ErrQueryDelegatorUnbondingDelegations = errors.Register(Codespace, 15, "query delegator unbonding delegations  error")
	ErrQueryRedelegations                 = errors.Register(Codespace, 16, "query redelegations  error")
	ErrQueryDelegatorValidators           = errors.Register(Codespace, 17, "query delegator validators  error")
	ErrQueryDelegatorValidator            = errors.Register(Codespace, 18, "query delegator validator error")
	ErrQueryHistoricalInfo                = errors.Register(Codespace, 19, "query historical info error")
	ErrQueryPool                          = errors.Register(Codespace, 20, "query pool error")
	ErrQueryParams                        = errors.Register(Codespace, 21, "query params  error")
	ErrGenConn                            = errors.Register(Codespace, 22, "generate conn error")
)
