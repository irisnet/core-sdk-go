package tx

import (
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.tx.v1beta1.Tx", (*sdk.Tx)(nil))
	registry.RegisterImplementations((*sdk.Tx)(nil), &Tx{})
}
