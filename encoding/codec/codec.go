package gongz

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cryptocodec "github.com/irisnet/core-sdk-go/crypto/codec"

	etherminttypes "github.com/irisnet/core-sdk-go/modules/ethermint/types"
)

// RegisterLegacyAminoCodec registers Interfaces from types, crypto, and SDK std.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {

	cryptocodec.RegisterCrypto(cdc)
	std.RegisterLegacyAminoCodec(cdc)

	// Register the auth modules msgs, requires import of x/auth/types.
	authtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	std.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)

	// Register the auth modules msgs, requires import of x/auth/types.
	authtypes.RegisterInterfaces(interfaceRegistry)
	// Register the ethermint account
	etherminttypes.RegisterInterfaces(interfaceRegistry)

}
