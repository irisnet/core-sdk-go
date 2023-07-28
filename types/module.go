package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
)

// The purpose of this interface is to convert the irishub system type to the user receiving type
// and standardize the user interface
type Response interface {
	Convert() interface{}
}

type SplitAble interface {
	Len() int
	Sub(begin, end int) SplitAble
}

type Module interface {
	Name() string
	RegisterInterfaceTypes(registry codectypes.InterfaceRegistry)
}

type KeyManager interface {
	Sign(name, password string, data []byte) ([]byte, cryptotypes.PubKey, error)
	Insert(name, password string) (string, string, error)
	Recover(name, password, mnemonic, hdPath string) (string, error)
	Import(name, password string, privKeyArmor string) (address string, err error)
	Export(name, password string) (privKeyArmor string, err error)
	Delete(name, password string) error
	Find(name, password string) (cryptotypes.PubKey, types.AccAddress, error)
	Add(name, password string) (address string, mnemonic string, err Error)
}
