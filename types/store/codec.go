package store

import (
	"github.com/irisnet/core-sdk-go/codec"
	cryptoAmino "github.com/irisnet/core-sdk-go/crypto/codec"
	"github.com/irisnet/core-sdk-go/crypto/hd"
	"github.com/tendermint/tendermint/crypto"
)

var cdc *codec.LegacyAmino

func init() {
	cdc = codec.NewLegacyAmino()
	cryptoAmino.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	cdc.Seal()
}

// RegisterCodec registers concrete types and interfaces on the given codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(hd.BIP44Params{}, "crypto/keys/hd/BIP44Params", nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = cdc.Unmarshal(pubKeyBytes, &pubKey)
	return
}

// encoding info
func marshalInfo(i Info) []byte {
	return cdc.MustMarshalLengthPrefixed(i)
}

// decoding info
func unmarshalInfo(bz []byte) (info Info, err error) {
	err = cdc.UnmarshalLengthPrefixed(bz, &info)
	return
}
