package nft_transfer

import (
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/common/codec"
	codectypes "github.com/irisnet/core-sdk-go/common/codec/types"
	cryptocodec "github.com/irisnet/core-sdk-go/common/crypto/codec"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
}

// RegisterInterfaces register the ibc nft-transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*types.Msg)(nil), &MsgTransfer{})
}
