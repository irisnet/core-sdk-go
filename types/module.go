package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// The purpose of this interface is to convert the irishub system type to the user receiving type
// and standardize the user interface
type Response interface {
	Convert() interface{}
}

type SplitAble interface {
	Len() int
	Sub(begin, end int) SplitAble
}

type Module interface {
	Name() string
	RegisterInterfaceTypes(registry codectypes.InterfaceRegistry)
}
