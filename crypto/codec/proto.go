package codec

import (
	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/crypto/keys/ed25519"
	"github.com/irisnet/core-sdk-go/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/crypto/keys/secp256r1"
	"github.com/irisnet/core-sdk-go/crypto/keys/sm2"
	cryptotypes "github.com/irisnet/core-sdk-go/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	var pk *cryptotypes.PubKey
	registry.RegisterInterface("cosmos.crypto.PubKey", pk)
	registry.RegisterImplementations(pk, &sm2.PubKey{})
	registry.RegisterImplementations(pk, &ed25519.PubKey{})
	registry.RegisterImplementations(pk, &secp256k1.PubKey{})
	secp256r1.RegisterInterfaces(registry)
}
