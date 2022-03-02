package feegrant

import (
	commoncodec "github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	commoncryptocodec "github.com/irisnet/core-sdk-go/common/crypto/codec"
	sdk "github.com/irisnet/core-sdk-go/types"
)

var (
	amino     = commoncodec.NewLegacyAmino()
	ModuleCdc = commoncodec.NewAminoCodec(amino)
)

func init() {
	commoncryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterInterfaces No duplicate registration
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgGrantAllowance{},
		&MsgRevokeAllowance{},
	)
}
