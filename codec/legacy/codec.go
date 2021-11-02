package legacy

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/core-sdk-go/codec"
	cryptocodec "github.com/irisnet/core-sdk-go/crypto/codec"
	cyptotypes "github.com/irisnet/core-sdk-go/crypto/types"
)

// Cdc defines a global generic sealed Amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *codec.LegacyAmino

func init() {
	Cdc = codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(Cdc)
	codec.RegisterEvidences(Cdc)
	Cdc.Seal()
}

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = Cdc.Unmarshal(privKeyBytes, &privKey)
	return
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey cyptotypes.PubKey, err error) {
	err = Cdc.Unmarshal(pubKeyBytes, &pubKey)
	return
}

func MarshalPubkey(pubkey crypto.PubKey) []byte {
	return Cdc.MustMarshal(pubkey)
}

func MarshalPrivKey(privKey crypto.PrivKey) []byte {
	return Cdc.MustMarshal(privKey)
}
