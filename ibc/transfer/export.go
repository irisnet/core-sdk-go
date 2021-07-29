package transfer

import (
	sdk "github.com/irisnet/core-sdk-go/types"
)

// expose transfer module api for user
type Client interface {
	sdk.Module

	CreatePool(request CreatePoolRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}

type CreatePoolRequest struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	LpTokenDenom   string    `json:"lp_token_denom"`
	StartHeight    int64     `json:"start_height"`
	RewardPerBlock sdk.Coins `json:"reward_per_block"`
	TotalReward    sdk.Coins `json:"total_reward"`
	Editable       bool      `json:"editable"`
	Creator        string    `json:"creator"`
}
