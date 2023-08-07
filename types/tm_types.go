package types

import (
	"encoding/hex"
	"strings"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	HexBytes      = tmbytes.HexBytes
	ABCIClient    = tmclient.ABCIClient
	SignClient    = tmclient.SignClient
	StatusClient  = tmclient.StatusClient
	NetworkClient = tmclient.NetworkClient
	Header        = tmtypes.Header
)

func MustHexBytesFrom(hexStr string) HexBytes {
	v, _ := hex.DecodeString(hexStr)
	return HexBytes(v)
}

func HexBytesFrom(hexStr string) (HexBytes, error) {
	v, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return HexBytes(v), nil
}

func HexStringFrom(bz []byte) string {
	return strings.ToUpper(hex.EncodeToString(bz))
}
