package gov

import (
	"github.com/irisnet/core-sdk-go/types/errors"
)

const Codespace = ModuleName

var (
	ErrTodo = errors.Register(Codespace, 2, "error todo")
)
