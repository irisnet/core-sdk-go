package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/irisnet/core-sdk-go/crypto/keys/ethsecp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmsm2 "github.com/tendermint/tendermint/crypto/sm2"
	tmsr25519 "github.com/tendermint/tendermint/crypto/sr25519"
)

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	//register tendermint public key
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)

	//register cosmos public key
	cdc.RegisterInterface((*cryptotypes.PubKey)(nil), nil)
	cdc.RegisterConcrete(tmed25519.PubKey{}, tmed25519.PubKeyName, nil)
	cdc.RegisterConcrete(tmsr25519.PubKey{}, tmsr25519.PubKeyName, nil)
	cdc.RegisterConcrete(tmsm2.PubKeySm2{}, tmsm2.PubKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PubKey{}, ed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PubKey{}, secp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&sm2.PubKey{}, sm2.PubKeyName, nil)
	cdc.RegisterConcrete(&ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&multisig.LegacyAminoPubKey{}, multisig.PubKeyAminoRoute, nil)

	//register private key
	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(tmed25519.PrivKey{}, tmed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PrivKey{}, ed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(tmsr25519.PrivKey{}, tmsr25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PrivKey{}, secp256k1.PrivKeyName, nil)
	cdc.RegisterConcrete(tmsm2.PrivKeySm2{}, tmsm2.PrivKeyName, nil)
	cdc.RegisterConcrete(&sm2.PrivKey{}, sm2.PrivKeyName, nil)
	cdc.RegisterConcrete(&ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName, nil)
}
