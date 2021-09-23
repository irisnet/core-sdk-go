package bank

import (
	"github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/modules/auth"
	sdk "github.com/irisnet/core-sdk-go/types"
)

// No duplicate registration
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSend{},
		&MsgMultiSend{},
	)

	registry.RegisterImplementations(
		(*auth.AccountI)(nil),
		&auth.BaseAccount{},
	)
}
