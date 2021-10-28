package bank

import (
	"github.com/irisnet/core-sdk-go/types/errors"
)

const Codespace = ModuleName

var (
	ErrQueryAccount     = errors.Register(Codespace, 1, "query account")
	ErrQueryTotalSupply = errors.Register(Codespace, 2, "query total supply option")
	ErrBuildAndSend     = errors.Register(Codespace, 3, "BuildAndSend error")
	ErrToMinCoin        = errors.Register(Codespace, 4, "ToMinCoin error")
	ErrGenConn          = errors.Register(Codespace, 5, "generate conn error")
)
