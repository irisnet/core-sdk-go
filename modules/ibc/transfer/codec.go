package transfer

import (
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/types"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*types.Msg)(nil),
		&MsgTransfer{},
	)
}
