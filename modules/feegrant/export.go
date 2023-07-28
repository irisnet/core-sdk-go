package feegrant

import (
	"github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	sdk "github.com/irisnet/core-sdk-go/types"
)

// Client expose fee grant module api for user
type Client interface {
	sdk.Module
	GrantAllowance(granter, grantee types.AccAddress, feeAllowance feegranttypes.FeeAllowanceI, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RevokeAllowance(granter, grantee types.AccAddress, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}
