package bank

import (
	"github.com/irisnet/core-sdk-go/types/errors"
)

const Codespace = ModuleName

var (
	ErrQueryAccount     = errors.Register(Codespace, 1, "query account")
	ErrQueryTotalSupply = errors.Register(Codespace, 1, "query total supply option")
	ErrBuildAndSend     = errors.Register(Codespace, 7, "BuildAndSend error")
	ErrToMinCoin        = errors.Register(Codespace, 22, "ToMinCoin error")
	ErrGenConn          = errors.Register(Codespace, 24, "generate conn error")
)
