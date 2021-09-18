package crypto

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto/armor"
)

const (
	// blockTypePrivKey = "TENDERMINT PRIVATE KEY"
	blockTypeKeyInfo = "TENDERMINT KEY INFO"
	blockTypePubKey  = "TENDERMINT PUBLIC KEY"

	defaultAlgo = "secp256k1"

	headerVersion = "version"
	headerType    = "type"
)

// BcryptSecurityParameter is security parameter var, and it can be changed within the lcd test.
// Making the bcrypt security parameter a var shouldn't be a security issue:
// One can't verify an invalid key by maliciously changing the bcrypt
// parameter during a runtime vulnerability. The main security
// threat this then exposes would be something that changes this during
// runtime before the user creates their key. This vulnerability must
// succeed to update this to that same value before every subsequent call
// to the keys command in future startups / or the attacker must get access
// to the filesystem. However, with a similar threat model (changing
// variables in runtime), one can cause the user to sign a different tx
// than what they see, which is a significantly cheaper attack then breaking
// a bcrypt hash. (Recall that the nonce still exists to break rainbow tables)
// For further notes on security parameter choice, see README.md
var BcryptSecurityParameter = 12

//-----------------------------------------------------------------
// add armor

// Armor the InfoBytes
func ArmorInfoBytes(bz []byte) string {
	header := map[string]string{
		headerType:    "Info",
		headerVersion: "0.0.0",
	}

	return armor.EncodeArmor(blockTypeKeyInfo, header, bz)
}

// Armor the PubKeyBytes
func ArmorPubKeyBytes(bz []byte, algo string) string {
	header := map[string]string{
		headerVersion: "0.0.1",
	}
	if algo != "" {
		header[headerType] = algo
	}

	return armor.EncodeArmor(blockTypePubKey, header, bz)
}

//-----------------------------------------------------------------
// remove armor

// Unarmor the InfoBytes
func UnarmorInfoBytes(armorStr string) ([]byte, error) {
	bz, header, err := unarmorBytes(armorStr, blockTypeKeyInfo)
	if err != nil {
		return nil, err
	}

	if header[headerVersion] != "0.0.0" {
		return nil, fmt.Errorf("unrecognized version: %v", header[headerVersion])
	}

	return bz, nil
}

// UnarmorPubKeyBytes returns the pubkey byte slice, a string of the algo type, and an error
func UnarmorPubKeyBytes(armorStr string) (bz []byte, algo string, err error) {
	bz, header, err := unarmorBytes(armorStr, blockTypePubKey)
	if err != nil {
		return nil, "", fmt.Errorf("couldn't unarmor bytes: %v", err)
	}

	switch header[headerVersion] {
	case "0.0.0":
		return bz, defaultAlgo, err
	case "0.0.1":
		if header[headerType] == "" {
			header[headerType] = defaultAlgo
		}

		return bz, header[headerType], err
	case "":
		return nil, "", fmt.Errorf("header's version field is empty")
	default:
		err = fmt.Errorf("unrecognized version: %v", header[headerVersion])
		return nil, "", err
	}
}

func unarmorBytes(armorStr, blockType string) (bz []byte, header map[string]string, err error) {
	bType, header, bz, err := armor.DecodeArmor(armorStr)
	if err != nil {
		return
	}

	if bType != blockType {
		err = fmt.Errorf("unrecognized armor type %q, expected: %q", bType, blockType)
		return
	}

	return
}
