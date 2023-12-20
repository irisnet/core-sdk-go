package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/core-sdk-go/crypto/keys/ethsecp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"
	tmsm2 "github.com/tendermint/tendermint/crypto/sm2"
)

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {

	cdc.RegisterConcrete(&ethsecp256k1.PubKey{},
		ethsecp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&ethsecp256k1.PrivKey{},
		ethsecp256k1.PrivKeyName, nil)

	cdc.RegisterConcrete(&sm2.PubKey{},
		sm2.PubKeyName, nil)
	cdc.RegisterConcrete(&sm2.PrivKey{},
		sm2.PrivKeyName, nil)

	cdc.RegisterConcrete(tmsm2.PubKeySm2{},
		tmsm2.PubKeyName, nil)
	cdc.RegisterConcrete(tmsm2.PrivKeySm2{},
		tmsm2.PrivKeyName, nil)
}
