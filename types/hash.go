package types

import (
	"crypto/sha256"

	"github.com/irisnet/core-sdk-go/types/errors"
	conv "github.com/irisnet/core-sdk-go/utils"
)

// Len is the length of base addresses
const Len = sha256.Size

// Addressable represents any type from which we can derive an address.
type Addressable interface {
	Address() []byte
}

// Hash creates a new address from address type and key
func Hash(typ string, key []byte) []byte {
	hasher := sha256.New()
	_, err := hasher.Write(conv.UnsafeStrToBytes(typ))
	// the error always nil, it's here only to satisfy the io.Writer interface
	errors.AssertNil(err)
	th := hasher.Sum(nil)

	hasher.Reset()
	_, err = hasher.Write(th)
	errors.AssertNil(err)
	_, err = hasher.Write(key)
	errors.AssertNil(err)
	return hasher.Sum(nil)
}

// Derive derives a new address from the main `address` and a derivation `key`.
func Derive(address []byte, key []byte) []byte {
	return Hash(conv.UnsafeBytesToStr(address), key)
}
