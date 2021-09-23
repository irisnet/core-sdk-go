package errors

import (
	"reflect"
)

const (
	SuccessCode              = 0
	internalCodespace        = UndefinedCodespace
	internalCode      uint32 = 1
)

type coder interface {
	Code() uint32
}

func Code(err error) uint32 {
	if errIsNil(err) {
		return SuccessCode
	}

	for {
		if c, ok := err.(coder); ok {
			return c.Code()
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return internalCode
		}
	}
}

type codespacer interface {
	Codespace() string
}

// abciCodespace tests if given error contains a codespace and returns the value of
// it if available. This function is testing for the causer interface as well
// and unwraps the error.
func Codespace(err error) string {
	if errIsNil(err) {
		return ""
	}

	for {
		if c, ok := err.(codespacer); ok {
			return c.Codespace()
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return internalCodespace
		}
	}
}

func errIsNil(err error) bool {
	if err == nil {
		return true
	}
	if val := reflect.ValueOf(err); val.Kind() == reflect.Ptr {
		return val.IsNil()
	}
	return false
}
