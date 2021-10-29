package types

import (
	proto "github.com/gogo/protobuf/proto"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)


// PubKey defines a public key and extends proto.Message.
type PubKey interface {
	proto.Message
	tmcrypto.PubKey
}

type PrivKey interface {
	proto.Message
	tmcrypto.PrivKey
}

type (
	Address = tmcrypto.Address
)
