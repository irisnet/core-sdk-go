package sm2

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"

	"github.com/tendermint/tendermint/crypto"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmsm2 "github.com/tendermint/tendermint/crypto/sm2"
	"github.com/tendermint/tendermint/crypto/tmhash"

	cryptotypes "github.com/irisnet/core-sdk-go/crypto/types"
)

const (
	PrivKeyName = "cosmos/PrivKeySm2"
	PubKeyName  = "cosmos/PubKeySm2"

	PrivKeySize   = 32
	PubKeySize    = 33
	SignatureSize = 64

	keyType = "sm2"
)

var _ cryptotypes.PrivKey = &PrivKey{}

// --------------------------------------------------------
func (privKey PrivKey) Type() string {
	return keyType
}

func (privKey PrivKey) Bytes() []byte {
	return privKey.Key
}

func (privKey PrivKey) Sign(msg []byte) ([]byte, error) {
	priv := privKey.GetPrivateKey()
	r, s, err := sm2.Sm2Sign(priv, msg, nil, rand.Reader)
	if err != nil {
		return nil, err
	}
	R := r.Bytes()
	S := s.Bytes()
	sig := make([]byte, 64)
	copy(sig[32-len(R):32], R[:])
	copy(sig[64-len(S):64], S[:])

	return sig, nil
}

func (privKey PrivKey) PubKey() tmcrypto.PubKey {
	priv := privKey.GetPrivateKey()
	compPubkey := sm2.Compress(&priv.PublicKey)

	pubkeyBytes := make([]byte, PubKeySize)
	copy(pubkeyBytes, compPubkey)

	return &PubKey{Key: pubkeyBytes}
}

func (privKey PrivKey) Equals(other tmcrypto.PrivKey) bool {
	if privKey.Type() != other.Type() {
		return false
	}

	return subtle.ConstantTimeCompare(privKey.Bytes(), other.Bytes()) == 1
}

func (privKey PrivKey) GetPrivateKey() *sm2.PrivateKey {
	k := new(big.Int).SetBytes(privKey.Key[:32])
	c := sm2.P256Sm2()
	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())

	return priv
}

func GenPrivKey() PrivKey {
	return genPrivKey(crypto.CReader())
}

func genPrivKey(rand io.Reader) PrivKey {
	seed := make([]byte, 32)
	if _, err := io.ReadFull(rand, seed); err != nil {
		panic(err)
	}

	privKey, err := sm2.GenerateKey(rand)
	if err != nil {
		panic(err)
	}

	privKeyBytes := make([]byte, PrivKeySize)
	copy(privKeyBytes, privKey.D.Bytes())

	return PrivKey{Key: privKeyBytes}
}

func GenPrivKeyFromSecret(secret []byte) PrivKey {
	one := new(big.Int).SetInt64(1)
	secHash := sha256.Sum256(secret)

	k := new(big.Int).SetBytes(secHash[:])
	n := new(big.Int).Sub(sm2.P256Sm2().Params().N, one)
	k.Mod(k, n)
	k.Add(k, one)

	return PrivKey{Key: k.Bytes()}
}

var _ cryptotypes.PubKey = &PubKey{}

// --------------------------------------------------------

func (pubKey PubKey) Address() crypto.Address {
	if len(pubKey.Key) != PubKeySize {
		panic("pubkey is incorrect size")
	}
	return crypto.Address(tmhash.SumTruncated(pubKey.Key))
}

func (pubKey PubKey) Bytes() []byte {
	return pubKey.Key
}

func (pubKey *PubKey) VerifySignature(msg []byte, sig []byte) bool {
	if len(sig) != SignatureSize {
		return false
	}

	publicKey := sm2.Decompress(pubKey.Key)
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])

	return sm2.Sm2Verify(publicKey, msg, nil, r, s)
}

func (pubKey PubKey) String() string {
	return fmt.Sprintf("PubKeySm2{%X}", pubKey.Key)
}

func (pubKey *PubKey) Type() string {
	return keyType
}

func (pubKey PubKey) Equals(other tmcrypto.PubKey) bool {
	if pubKey.Type() != other.Type() {
		return false
	}

	return subtle.ConstantTimeCompare(pubKey.Bytes(), other.Bytes()) == 1
}

// AsTmPubKey converts our own PubKey into a Tendermint ED25519 pubkey.
func (pubKey *PubKey) AsTmPubKey() crypto.PubKey {
	var pubkey tmsm2.PubKeySm2
	copy(pubkey[:], pubKey.Key)
	return pubkey
}
