package feegrant

import (
	sdk "github.com/irisnet/core-sdk-go/types"
)

// Client expose fee grant module api for user
type Client interface {
	sdk.Module
	GrantAllowance(granter, grantee sdk.AccAddress,feeAllowance FeeAllowanceI, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RevokeAllowance(granter, grantee sdk.AccAddress, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}
