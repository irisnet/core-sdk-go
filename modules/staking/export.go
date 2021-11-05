package staking

import (
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/core-sdk-go/types"
)

// expose Staking module api for user
type Client interface {
	types.Module

	CreateValidator(request CreateValidatorRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)
	EditValidator(request EditValidatorRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)
	Delegate(request DelegateRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)
	Undelegate(request UndelegateRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)
	BeginRedelegate(request BeginRedelegateRequest, baseTx types.BaseTx) (ctypes.ResultTx, error)

	QueryValidators(status string, page, size uint64) (QueryValidatorsResp, error)
	QueryValidator(validatorAddr string) (QueryValidatorResp, error)
	QueryValidatorDelegations(validatorAddr string, page, size uint64) (QueryValidatorDelegationsResp, error)
	QueryValidatorUnbondingDelegations(validatorAddr string, page, size uint64) (QueryValidatorUnbondingDelegationsResp, error)
	QueryDelegation(delegatorAddr string, validatorAddr string) (QueryDelegationResp, error)
	QueryUnbondingDelegation(delegatorAddr string, validatorAddr string) (QueryUnbondingDelegationResp, error)
	QueryDelegatorDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorDelegationsResp, error)
	QueryDelegatorUnbondingDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorUnbondingDelegationsResp, error)
	QueryRedelegations(request QueryRedelegationsReq) (QueryRedelegationsResp, error)
	QueryDelegatorValidators(delegatorAddr string, page, size uint64) (QueryDelegatorValidatorsResp, error)
	QueryDelegatorValidator(delegatorAddr string, validatorAddr string) (QueryValidatorResp, error)
	QueryHistoricalInfo(height int64) (QueryHistoricalInfoResp, error)
	QueryPool() (QueryPoolResp, error)
	QueryParams() (QueryParamsResp, error)
}

type CreateValidatorRequest struct {
	Moniker           string        `json:"moniker"`
	Rate              types.Dec     `json:"rate"`
	MaxRate           types.Dec     `json:"max_rate"`
	MaxChangeRate     types.Dec     `json:"max_change_rate"`
	MinSelfDelegation types.Int     `json:"min_self_delegation"`
	Pubkey            string        `json:"pubkey"`
	Value             types.DecCoin `json:"value"`
}

type EditValidatorRequest struct {
	Moniker           string    `json:"moniker"`
	Identity          string    `json:"identity"`
	Website           string    `json:"website"`
	SecurityContact   string    `json:"security_contact"`
	Details           string    `json:"details"`
	CommissionRate    types.Dec `json:"commission_rate"`
	MinSelfDelegation types.Int `json:"min_self_delegation"`
}

type DelegateRequest struct {
	ValidatorAddr string        `json:"validator_address"`
	Amount        types.DecCoin `json:"amount"`
}

type UndelegateRequest struct {
	ValidatorAddr string        `json:"validator_address"`
	Amount        types.DecCoin `json:"amount"`
}

type BeginRedelegateRequest struct {
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Amount              types.DecCoin
}

type (
	description struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	}
	commission struct {
		commissionRates
		UpdateTime time.Time `json:"update_time"`
	}
	commissionRates struct {
		Rate          types.Dec `json:"rate"`
		MaxRate       types.Dec `json:"max_rate"`
		MaxChangeRate types.Dec `json:"max_change_rate"`
	}

	QueryValidatorsResp struct {
		Validators []QueryValidatorResp `json:"validators"`
		Total      uint64               `json:"total"`
	}

	QueryValidatorResp struct {
		OperatorAddress   string      `json:"operator_address"`
		ConsensusPubkey   string      `json:"consensus_pubkey"`
		Jailed            bool        `json:"jailed"`
		Status            string      `json:"status"`
		Tokens            types.Int   `json:"tokens"`
		DelegatorShares   types.Dec   `json:"delegator_shares"`
		Description       description `json:"description"`
		UnbondingHeight   int64       `json:"unbonding_height"`
		UnbondingTime     time.Time   `json:"unbonding_time"`
		Commission        commission  `json:"commission"`
		MinSelfDelegation types.Int   `json:"min_self_delegation"`
	}
)

type (
	delegation struct {
		DelegatorAddress string    `json:"delegator_address"`
		Shares           types.Dec `json:"shares"`
		ValidatorAddress string    `json:"validator_address"`
	}

	QueryDelegationResp struct {
		Delegation delegation `json:"delegation"`
		Balance    types.Coin `json:"balance"`
	}

	QueryValidatorDelegationsResp struct {
		DelegationResponses []QueryDelegationResp `json:"delegation_responses"`
		Total               uint64                `json:"total"`
	}
)

type (
	unbondingDelegationEntry struct {
		CreationHeight int64     `json:"creation_height"`
		CompletionTime time.Time `json:"completion_time"`
		InitialBalance types.Int `json:"initial_balance"`
		Balance        types.Int `json:"balance"`
	}

	QueryUnbondingDelegationResp struct {
		DelegatorAddress string                     `json:"delegator_address"`
		ValidatorAddress string                     `json:"validator_address"`
		Entries          []unbondingDelegationEntry `json:"entries"`
	}

	QueryValidatorUnbondingDelegationsResp struct {
		UnbondingResponses []QueryUnbondingDelegationResp `json:"unbonding_responses"`
		Total              uint64                         `json:"total"`
	}
)

type QueryDelegatorDelegationsResp struct {
	DelegationResponses []QueryDelegationResp `json:"delegation_responses"`
	Total               uint64                `json:"total"`
}

type QueryDelegatorUnbondingDelegationsResp struct {
	UnbondingDelegations []QueryUnbondingDelegationResp `json:"unbonding_delegations"`
	Total                uint64                         `json:"total"`
}

type (
	QueryRedelegationsReq struct {
		DelegatorAddr    string `json:"delegator_addr"`
		SrcValidatorAddr string `json:"src_validator_addr"`
		DstValidatorAddr string `json:"dst_validator_addr"`
		Page             uint64 `json:"page"`
		Size             uint64 `json:"size"`
	}

	QueryRedelegationsResp struct {
		RedelegationResponses []RedelegationResp `json:"redelegation_responses"`
		Total                 uint64             `json:"total"`
	}

	redelegationEntry struct {
		CreationHeight int64     `json:"creation_height"`
		CompletionTime time.Time `json:"completion_time"`
		InitialBalance types.Int `json:"initial_balance"`
		SharesDst      types.Dec `json:"shares_dst"`
	}
	redelegationEntryResponse struct {
		RedelegationEntry redelegationEntry `json:"redelegation_entry"`
		Balance           types.Int         `json:"balance"`
	}
	redelegation struct {
		DelegatorAddress    string              `json:"delegator_address"`
		ValidatorSrcAddress string              `json:"validator_src_address"`
		ValidatorDstAddress string              `json:"validator_dst_address"`
		Entries             []redelegationEntry `json:"entries"`
	}

	RedelegationResp struct {
		Redelegation redelegation                `json:"redelegation"`
		Entries      []redelegationEntryResponse `json:"entries"`
	}
)

type QueryDelegatorValidatorsResp struct {
	Validator []QueryValidatorResp `json:"validator"`
	Total     uint64               `json:"total"`
}

type QueryHistoricalInfoResp struct {
	Header types.Header         `json:"header"`
	Valset []QueryValidatorResp `json:"valset"`
}

type QueryPoolResp struct {
	NotBondedTokens types.Int `json:"not_bonded_tokens"`
	BondedTokens    types.Int `json:"bonded_tokens"`
}

type QueryParamsResp struct {
	UnbondingTime     time.Duration `json:"unbonding_time"`
	MaxValidators     uint32        `json:"max_validators"`
	MaxEntries        uint32        `json:"max_entries"`
	HistoricalEntries uint32        `json:"historical_entries"`
	BondDenom         string        `json:"bond_denom"`
}
