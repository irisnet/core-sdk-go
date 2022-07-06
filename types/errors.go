package types

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	// RootCodespace is the codespace for all errors defined in irita
	RootCodespace = "sdk"

	OK                      Code = 0
	Internal                Code = 1
	TxDecode                Code = 2
	InvalidSequence         Code = 3
	Unauthorized            Code = 4
	InsufficientFunds       Code = 5
	UnknownRequest          Code = 6
	InvalidAddress          Code = 7
	InvalidPubkey           Code = 8
	UnknownAddress          Code = 9
	InvalidCoins            Code = 10
	OutOfGas                Code = 11
	MemoTooLarge            Code = 12
	InsufficientFee         Code = 13
	TooManySignatures       Code = 14
	NoSignatures            Code = 15
	ErrJsonMarshal          Code = 16
	ErrJsonUnmarshal        Code = 17
	InvalidRequest          Code = 18
	TxInMempoolCache        Code = 19
	MempoolIsFull           Code = 20
	TxTooLarge              Code = 21
	KeyNotFound             Code = 22
	WrongPassword           Code = 23
	InvalidSigner           Code = 24
	InvalidGasAdjustment    Code = 25
	InvalidHeight           Code = 26
	InvalidVersion          Code = 27
	InvalidChainID          Code = 28
	InvalidType             Code = 29
	TxTimeoutHeight         Code = 30
	UnknownExtensionOptions Code = 31
	WrongSequence           Code = 32
	PackAny                 Code = 33
	UnpackAny               Code = 34
	Logic                   Code = 35
	Conflict                Code = 36
	NotSupported            Code = 37
	NotFound                Code = 38
	IO                      Code = 39
	AppConfig               Code = 40
	Panic                   Code = 111222
)

var (
	// errUnknown = register(RootCodespace, 111222, "unknown error")
	errInvalid  = register(RootCodespace, 999999, "sdk check error")
	wrongSeqMsg = "account sequence mismatch, expected"
)

var (
	Success                    = register(RootCodespace, OK, "success")
	ErrInternal                = register(RootCodespace, Internal, "internal")
	ErrTxDecode                = register(RootCodespace, TxDecode, "tx parse error")
	ErrInvalidSequence         = register(RootCodespace, InvalidSequence, "invalid sequence")
	ErrUnauthorized            = register(RootCodespace, Unauthorized, "unauthorized")
	ErrInsufficientFunds       = register(RootCodespace, InsufficientFunds, "insufficient funds")
	ErrUnknownRequest          = register(RootCodespace, UnknownRequest, "unknown request")
	ErrInvalidAddress          = register(RootCodespace, InvalidAddress, "invalid address")
	ErrInvalidPubKey           = register(RootCodespace, InvalidPubkey, "invalid pubkey")
	ErrUnknownAddress          = register(RootCodespace, UnknownAddress, "unknown address")
	ErrInvalidCoins            = register(RootCodespace, InvalidCoins, "invalid coins")
	ErrOutOfGas                = register(RootCodespace, OutOfGas, "out of gas")
	ErrMemoTooLarge            = register(RootCodespace, MemoTooLarge, "memo too large")
	ErrInsufficientFee         = register(RootCodespace, InsufficientFee, "insufficient fee")
	ErrTooManySignatures       = register(RootCodespace, TooManySignatures, "maximum number of signatures exceeded")
	ErrNoSignatures            = register(RootCodespace, NoSignatures, "no signatures supplied")
	ErrJSONMarshal             = register(RootCodespace, ErrJsonMarshal, "failed to marshal JSON bytes")
	ErrJSONUnmarshal           = register(RootCodespace, ErrJsonUnmarshal, "failed to unmarshal JSON bytes")
	ErrInvalidRequest          = register(RootCodespace, InvalidRequest, "invalid request")
	ErrTxInMempoolCache        = register(RootCodespace, TxInMempoolCache, "tx already in mempool")
	ErrMempoolIsFull           = register(RootCodespace, MempoolIsFull, "mempool is full")
	ErrTxTooLarge              = register(RootCodespace, TxTooLarge, "tx too large")
	ErrKeyNotFound             = register(RootCodespace, KeyNotFound, "key not found")
	ErrWrongPassword           = register(RootCodespace, WrongPassword, "invalid account password")
	ErrorInvalidSigner         = register(RootCodespace, InvalidSigner, "tx intended signer does not match the given signer")
	ErrorInvalidGasAdjustment  = register(RootCodespace, InvalidGasAdjustment, "invalid gas adjustment")
	ErrInvalidHeight           = register(RootCodespace, InvalidHeight, "invalid height")
	ErrInvalidVersion          = register(RootCodespace, InvalidVersion, "invalid version")
	ErrInvalidChainID          = register(RootCodespace, InvalidChainID, "invalid chain-id")
	ErrInvalidType             = register(RootCodespace, InvalidType, "invalid type")
	ErrTxTimeoutHeight         = register(RootCodespace, TxTimeoutHeight, "tx timeout height")
	ErrUnknownExtensionOptions = register(RootCodespace, UnknownExtensionOptions, "unknown extension options")
	ErrWrongSequence           = register(RootCodespace, WrongSequence, "incorrect account sequence")
	ErrPackAny                 = register(RootCodespace, PackAny, "failed packing protobuf message to Any")
	ErrUnpackAny               = register(RootCodespace, UnpackAny, "failed unpacking protobuf message from Any")
	ErrLogic                   = register(RootCodespace, Logic, "internal logic error")
	ErrConflict                = register(RootCodespace, Conflict, "conflict")
	ErrNotSupported            = register(RootCodespace, NotSupported, "feature not supported")
	ErrNotFound                = register(RootCodespace, NotFound, "not found")
	ErrIO                      = register(RootCodespace, IO, "Internal IO error")
	ErrAppConfig               = register(RootCodespace, AppConfig, "error in app.toml")
	ErrPanic                   = register(RootCodespace, Panic, "panic")
)

type Code uint32

// Error represents a root error.
//
// Weave framework is using root error to categorize issues. Each instance
// created during the runtime should wrap one of the declared root errors. This
// allows error tests and returning all errors to the client in a safe manner.
//
// All popular root errors are declared in this package. If an extension has to
// declare a custom root error, always use register function to ensure
// error code uniqueness.
type Error interface {
	Error() string
	Code() uint32
	Codespace() string
	WrapfError(string) Error
}

// GetError is used to covert irita error to sdk error
func GetError(codespace string, code uint32, log ...string) Error {
	errID := errorID(codespace, code)
	err, ok := usedCodes[errID]
	if !ok {
		return Wrap(errors.New(strings.Join(log, ",")))
	}

	return err.WrapfError(strings.Join(log, "."))
}

// Wrap extends given error with an additional information.
//
// If the wrapped error does not provide ABCICode method (ie. stdlib errors),
// it will be labeled as internal error.
//
// If err is nil, this returns nil, avoiding the need for an if statement when
// wrapping a error returned at the end of a function
func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	code := errInvalid.Code()
	codespace := errInvalid.Codespace()

	if strings.Contains(err.Error(), wrongSeqMsg) {
		return sdkError{
			codespace: RootCodespace,
			code:      uint32(WrongSequence),
			desc:      err.Error(),
		}
	}

	e, ok := err.(sdkError)
	if ok {
		code = e.Code()
		codespace = e.Codespace()
	}

	return sdkError{
		code:      code,
		codespace: codespace,
		desc:      err.Error(),
	}
}

func WrapWithMessage(err error, format string, args ...interface{}) Error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(errors.WithMessage(err, desc))
}

// Wrapf extends given error with an additional information.
//
// This function works like Wrap function with additional functionality of
// formatting the input as specified.
func Wrapf(format string, args ...interface{}) Error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(errors.New(desc))
}

type sdkError struct {
	codespace string
	code      uint32
	desc      string
}

func (e sdkError) WrapfError(desc string) Error {
	e.desc = fmt.Sprintf("%s: %s", e.desc, desc)
	return e
}

func (e sdkError) Error() string {
	return e.desc
}

func (e sdkError) Code() uint32 {
	return e.code
}

func (e sdkError) Codespace() string {
	return e.codespace
}

// register returns an error instance that should be used as the base for
// creating error instances during runtime.
//
// Popular root errors are declared in this package, but extensions may want to
// declare custom codes. This function ensures that no error code is used
// twice. Attempt to reuse an error code results in panic.
//
// Use this function only during a program startup phase.
func register(codespace string, code Code, description string) Error {
	err := sdkError{
		codespace: codespace,
		code:      uint32(code),
		desc:      description,
	}
	setUsed(err)

	return err
}

// usedCodes is keeping track of used codes to ensure their uniqueness. No two
// error instances should share the same (codespace, code) tuple.
var usedCodes = map[string]Error{}

func errorID(codespace string, code uint32) string {
	return fmt.Sprintf("%s:%d", codespace, code)
}

func setUsed(err Error) {
	usedCodes[errorID(err.Codespace(), err.Code())] = err
}

func CatchPanic(fn func(errMsg string)) {
	if err := recover(); err != nil {
		var msg string
		switch e := err.(type) {
		case error:
			msg = e.Error()
		case string:
			msg = e
		}
		fn(msg)
	}
}

func RegisterErr(codespace string, code Code, description string) Error {
	return register(codespace, code, description)
}
