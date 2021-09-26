package tx

import (
	sdk "github.com/irisnet/core-sdk-go/types"
	signingtypes "github.com/irisnet/core-sdk-go/types/tx/signing"
)

// signModeDirectHandler defines the SIGN_MODE_DIRECT SignModeHandler
type signModeDirectHandler struct{}

var _ sdk.SignModeHandler = signModeDirectHandler{}

// DefaultMode implements SignModeHandler.DefaultMode
func (signModeDirectHandler) DefaultMode() signingtypes.SignMode {
	return signingtypes.SignMode_SIGN_MODE_DIRECT
}

// Modes implements SignModeHandler.Modes
func (signModeDirectHandler) Modes() []signingtypes.SignMode {
	return []signingtypes.SignMode{signingtypes.SignMode_SIGN_MODE_DIRECT}
}
