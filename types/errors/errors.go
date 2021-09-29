package errors

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// RootCodespace is the codespace for all errors defined in this package
const RootCodespace = "sdk"

// UndefinedCodespace when we explicitly declare no codespace
const UndefinedCodespace = "undefined"

var (
	ErrPanic           = Register(UndefinedCodespace, 111222, "panic")
	ErrTodo            = Register(RootCodespace, 100, "error todo")
	ErrInvalidAddress  = Register(RootCodespace, 7, "invalid address")
	ErrInvalidRequest  = Register(RootCodespace, 18, "invalid request")
	ErrPackAny         = Register(RootCodespace, 33, "failed packing protobuf message to Any")
	ErrUnpackAny       = Register(RootCodespace, 34, "failed unpacking protobuf message from Any")
	ErrInvalidPubKey   = Register(RootCodespace, 8, "invalid pubkey")
	ErrInvalidType     = Register(RootCodespace, 29, "invalid type")
	ErrInvalidSequence = Register(RootCodespace, 3, "invalid sequence")
	ErrTxTooLarge      = Register(RootCodespace, 21, "tx too large")
	ErrWrongSequence   = Register(RootCodespace, 32, "incorrect account sequence")
	ErrWrongPassword   = Register(RootCodespace, 23, "invalid account password")

	// ErrTxDecode = Register(RootCodespace, 2, "tx parse error")
	// ErrUnauthorized = Register(RootCodespace, 4, "unauthorized")
	// ErrInsufficientFunds = Register(RootCodespace, 5, "insufficient funds")
	// ErrUnknownRequest = Register(RootCodespace, 6, "unknown request")
	// ErrUnknownAddress = Register(RootCodespace, 9, "unknown address")
	// ErrInvalidCoins = Register(RootCodespace, 10, "invalid coins")
	// ErrOutOfGas = Register(RootCodespace, 11, "out of gas")
	// ErrMemoTooLarge = Register(RootCodespace, 12, "memo too large")
	// ErrInsufficientFee = Register(RootCodespace, 13, "insufficient fee")
	// ErrTooManySignatures = Register(RootCodespace, 14, "maximum number of signatures exceeded")
	// ErrNoSignatures = Register(RootCodespace, 15, "no signatures supplied")
	// ErrJSONMarshal = Register(RootCodespace, 16, "failed to marshal JSON bytes")
	// ErrJSONUnmarshal = Register(RootCodespace, 17, "failed to unmarshal JSON bytes")
	// ErrTxInMempoolCache = Register(RootCodespace, 19, "tx already in mempool")
	// ErrMempoolIsFull = Register(RootCodespace, 20, "mempool is full")
	// ErrKeyNotFound = Register(RootCodespace, 22, "key not found")
	// ErrorInvalidSigner = Register(RootCodespace, 24, "tx intended signer does not match the given signer")
	// ErrorInvalidGasAdjustment = Register(RootCodespace, 25, "invalid gas adjustment")
	// ErrInvalidHeight = Register(RootCodespace, 26, "invalid height")
	// ErrInvalidVersion = Register(RootCodespace, 27, "invalid version")
	// ErrInvalidChainID = Register(RootCodespace, 28, "invalid chain-id")
	// ErrTxTimeoutHeight = Register(RootCodespace, 30, "tx timeout height")
	// ErrUnknownExtensionOptions = Register(RootCodespace, 31, "unknown extension options")
	// ErrLogic = Register(RootCodespace, 35, "internal logic error")
	// ErrConflict = Register(RootCodespace, 36, "conflict")
	// ErrNotSupported = Register(RootCodespace, 37, "feature not supported")
	// ErrNotFound = Register(RootCodespace, 38, "not found")
	// ErrIO = Register(RootCodespace, 39, "Internal IO error")

)

// Register returns an error instance that should be used as the base for
// creating error instances during runtime.
//
// Popular root errors are declared in this package, but extensions may want to
// declare custom codes. This function ensures that no error code is used
// twice. Attempt to reuse an error code results in panic.
//
// Use this function only during a program startup phase.
func Register(codespace string, code uint32, description string) *Error {
	// TODO - uniqueness is (codespace, code) combo
	if e := getUsed(codespace, code); e != nil {
		panic(fmt.Sprintf("error with code %d is already registered: %q", code, e.desc))
	}

	err := New(codespace, code, description)
	setUsed(err)

	return err
}

// usedCodes is keeping track of used codes to ensure their uniqueness. No two
// error instances should share the same (codespace, code) tuple.
var usedCodes = map[string]*Error{}

func errorID(codespace string, code uint32) string {
	return fmt.Sprintf("%s:%d", codespace, code)
}

func getUsed(codespace string, code uint32) *Error {
	return usedCodes[errorID(codespace, code)]
}

func setUsed(err *Error) {
	usedCodes[errorID(err.codespace, err.code)] = err
}

// Error represents a root error.
//
// Weave framework is using root error to categorize issues. Each instance
// created during the runtime should wrap one of the declared root errors. This
// allows error tests and returning all errors to the client in a safe manner.
//
// All popular root errors are declared in this package. If an extension has to
// declare a custom root error, always use Register function to ensure
// error code uniqueness.
type Error struct {
	codespace string
	code      uint32
	desc      string
}

func New(codespace string, code uint32, desc string) *Error {
	return &Error{codespace: codespace, code: code, desc: desc}
}

func (e Error) Error() string {
	return e.desc
}

func (e Error) Code() uint32 {
	return e.code
}

func (e Error) Codespace() string {
	return e.codespace
}

// Is check if given error instance is of a given kind/type. This involves
// unwrapping given error using the Cause method if available.
func (e *Error) Is(err error) bool {
	// Reflect usage is necessary to correctly compare with
	// a nil implementation of an error.
	if e == nil {
		return isNilErr(err)
	}

	for {
		if err == e {
			return true
		}

		// If this is a collection of errors, this function must return
		// true if at least one from the group match.
		if u, ok := err.(unpacker); ok {
			for _, er := range u.Unpack() {
				if e.Is(er) {
					return true
				}
			}
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return false
		}
	}
}

// Wrap extends this error with an additional information.
// It's a handy function to call Wrap with sdk errors.
func (e Error) Wrap(desc string) error { return Wrap(e, desc) }

// Wrapf extends this error with an additional information.
// It's a handy function to call Wrapf with sdk errors.
func (e Error) Wrapf(desc string, args ...interface{}) error { return Wrapf(e, desc, args...) }

func isNilErr(err error) bool {
	// Reflect usage is necessary to correctly compare with
	// a nil implementation of an error.
	if err == nil {
		return true
	}
	if reflect.ValueOf(err).Kind() == reflect.Struct {
		return false
	}
	return reflect.ValueOf(err).IsNil()
}

// Wrap extends given error with an additional information.
//
// If the wrapped error does not provide ABCICode method (ie. stdlib errors),
// it will be labeled as internal error.
//
// If err is nil, this returns nil, avoiding the need for an if statement when
// wrapping a error returned at the end of a function
func Wrap(err error, description string) error {
	if err == nil {
		return nil
	}

	// If this error does not carry the stacktrace information yet, attach
	// one. This should be done only once per error at the lowest frame
	// possible (most inner wrap).
	if stackTrace(err) == nil {
		err = errors.WithStack(err)
	}

	return &wrappedError{
		parent: err,
		msg:    description,
	}
}

// Wrapf extends given error with an additional information.
//
// This function works like Wrap function with additional functionality of
// formatting the input as specified.
func Wrapf(err error, format string, args ...interface{}) error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(err, desc)
}

type wrappedError struct {
	// This error layer description.
	msg string
	// The underlying error that triggered this one.
	parent error
}

func (e *wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.parent.Error())
}

func (e *wrappedError) Cause() error {
	return e.parent
}

// Is reports whether any error in e's chain matches a target.
func (e *wrappedError) Is(target error) bool {
	if e == target {
		return true
	}

	w := e.Cause()
	for {
		if w == target {
			return true
		}

		x, ok := w.(causer)
		if ok {
			w = x.Cause()
		}
		if x == nil {
			return false
		}
	}
}

// Unwrap implements the built-in errors.Unwrap
func (e *wrappedError) Unwrap() error {
	return e.parent
}

// Recover captures a panic and stop its propagation. If panic happens it is
// transformed into a ErrPanic instance and assigned to given error. Call this
// function using defer in order to work as expected.
func Recover(err *error) {
	if r := recover(); r != nil {
		*err = Wrapf(ErrPanic, "%v", r)
	}
}

// WithType is a helper to augment an error with a corresponding type message
func WithType(err error, obj interface{}) error {
	return Wrap(err, fmt.Sprintf("%T", obj))
}

// IsOf checks if a received error is caused by one of the target errors.
// It extends the errors.Is functionality to a list of errors.
func IsOf(received error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(received, t) {
			return true
		}
	}
	return false
}

// causer is an interface implemented by an error that supports wrapping. Use
// it to test if an error wraps another error instance.
type causer interface {
	Cause() error
}

type unpacker interface {
	Unpack() []error
}

// stackTrace returns the first found stack trace frame carried by given error
// or any wrapped error. It returns nil if no stack trace is found.
func stackTrace(err error) errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	for {
		if st, ok := err.(stackTracer); ok {
			return st.StackTrace()
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return nil
		}
	}
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
