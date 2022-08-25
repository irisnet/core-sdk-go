package codec

import (
	"github.com/irisnet/core-sdk-go/common/codec"
	commoncrypto "github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/ed25519"
	ethsecp256k1 "github.com/irisnet/core-sdk-go/common/crypto/keys/eth_secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/multisig"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/sm2"
	cryptotypes "github.com/irisnet/core-sdk-go/common/crypto/types"
	"github.com/tendermint/tendermint/crypto"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	tmsm2 "github.com/tendermint/tendermint/crypto/sm2"
	tmsr25519 "github.com/tendermint/tendermint/crypto/sr25519"
)

var amino *commoncrypto.LegacyAmino

func init() {
	amino = commoncrypto.NewLegacyAmino()
	RegisterCrypto(amino)
}

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

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = amino.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = amino.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}

func MarshalPubkey(pubkey crypto.PubKey) []byte {
	return amino.MustMarshalBinaryBare(pubkey)
}

func MarshalPrivKey(privKey crypto.PrivKey) []byte {
	return amino.MustMarshalBinaryBare(privKey)
}
