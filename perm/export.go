package perm

import sdk "github.com/irisnet/core-sdk-go/types"

// Client expose fee grant module api for user
type Client interface {
	sdk.Module
	AssignRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}
