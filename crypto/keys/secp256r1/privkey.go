package secp256r1

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/core-sdk-go/crypto/keys/internal/ecdsa"
)

// GenPrivKey generates a new secp256r1 private key. It uses operating system randomness.
func GenPrivKey() (*PrivKey, error) {
	key, err := ecdsa.GenPrivKey(secp256r1)
	return &PrivKey{&ecdsaSK{key}}, err
}

// PubKey implements SDK PrivKey interface.
func (m *PrivKey) PubKey() tmcrypto.PubKey {
	return &PubKey{&ecdsaPK{m.Secret.PubKey()}}
}

// String implements SDK proto.Message interface.
func (m *PrivKey) String() string {
	return m.Secret.String(KeyType)
}

// Type returns key type KeyType. Implements SDK PrivKey interface.
func (m *PrivKey) Type() string {
	return KeyType
}

// Sign hashes and signs the message usign ECDSA. Implements sdk.PrivKey interface.
func (m *PrivKey) Sign(msg []byte) ([]byte, error) {
	return m.Secret.Sign(msg)
}

// Bytes serialize the private key.
func (m *PrivKey) Bytes() []byte {
	if m == nil {
		return nil
	}
	return m.Secret.Bytes()
}

// Equals implements SDK PrivKey interface.
func (m *PrivKey) Equals(other tmcrypto.PrivKey) bool {
	sk2, ok := other.(*PrivKey)
	if !ok {
		return false
	}
	return m.Secret.Equal(&sk2.Secret.PrivateKey)
}

type ecdsaSK struct {
	ecdsa.PrivKey
}

// Size implements proto.Marshaler interface
func (sk *ecdsaSK) Size() int {
	if sk == nil {
		return 0
	}
	return fieldSize
}

// Unmarshal implements proto.Marshaler interface
func (sk *ecdsaSK) Unmarshal(bz []byte) error {
	return sk.PrivKey.Unmarshal(bz, secp256r1, fieldSize)
}
