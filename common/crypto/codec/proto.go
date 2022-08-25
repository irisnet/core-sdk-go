package codec

import (
	ethsecp256k1 "github.com/irisnet/core-sdk-go/common/crypto/keys/eth_secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/sm2"
	tmcrypto "github.com/tendermint/tendermint/crypto"

	codectypes "github.com/irisnet/core-sdk-go/common/codec/types"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/ed25519"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/multisig"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/secp256k1"
	cryptotypes "github.com/irisnet/core-sdk-go/common/crypto/types"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// TODO We now register both Tendermint's PubKey and our own PubKey. In the
	// long-term, we should move away from Tendermint's PubKey, and delete
	// these lines.
	registry.RegisterInterface("tendermint.crypto.Pubkey", (*tmcrypto.PubKey)(nil))
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &sm2.PubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &multisig.LegacyAminoPubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &ethsecp256k1.PubKey{})

	registry.RegisterInterface("cosmos.crypto.Pubkey", (*cryptotypes.PubKey)(nil))
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &sm2.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &multisig.LegacyAminoPubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ethsecp256k1.PubKey{})

}
