package keyring_test

import (
	"testing"

	"github.com/irisnet/core-sdk-go/common/address/irishub"

	"github.com/irisnet/core-sdk-go/crypto/keyring"

	"github.com/stretchr/testify/assert"
)

func TestNewMnemonicKeyManager(t *testing.T) {
	mnemonic := "nerve leader thank marriage spice task van start piece crowd run hospital control outside cousin romance left choice poet wagon rude climb leisure spring"

	km, err := keyring.NewMnemonicKeyManager(mnemonic, "secp256k1")
	assert.NoError(t, err)

	pubKey := km.ExportPubKey()

	address := irishub.Bech32Address(pubKey.Address().Bytes())

	assert.Equal(t, "iaa1y9kd9uy7a4qnjp0z5yjx5jhrkv2ycdkzqc0h8z", address)
}
