package codec

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/sr25519"

	"github.com/irisnet/core-sdk-go/codec"
	"github.com/irisnet/core-sdk-go/crypto/keys/ed25519"
	"github.com/irisnet/core-sdk-go/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"
	cryptotypes "github.com/irisnet/core-sdk-go/crypto/types"
)

// RegisterCrypto registers all crypto dependency types with the provided Amino codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*tmcrypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*cryptotypes.PubKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PubKey{}, sr25519.PubKeyName, nil)
	cdc.RegisterConcrete(tmed25519.PubKey{}, tmed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PubKey{}, ed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PubKey{}, secp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&sm2.PubKey{}, sm2.PubKeyName, nil)

	cdc.RegisterInterface((*tmcrypto.PrivKey)(nil), nil)
	cdc.RegisterInterface((*cryptotypes.PrivKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PrivKey{}, sr25519.PrivKeyName, nil)
	cdc.RegisterConcrete(tmed25519.PrivKey{}, tmed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PrivKey{}, ed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PrivKey{}, secp256k1.PrivKeyName, nil)
	cdc.RegisterConcrete(&sm2.PrivKey{}, sm2.PrivKeyName, nil)
}
