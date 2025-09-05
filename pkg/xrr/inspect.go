package xrr

import (
	"errors"
)

// GetCode returns error code. If the error implements [Coder] interface, it
// will be used. Otherwise, it goes through error chains (and joins) to find
// the first non-empty error code.
//
// For nil error it will return an empty string. If the error code is not found,
// returns [ECGeneric] error code.
func GetCode(err error) string {
	if err == nil || isNil(err) {
		return ""
	}

	switch e := err.(type) { // nolint: errorlint
	case Coder:
		return e.ErrorCode()

	case interface{ Unwrap() []error }:
		// Get the first code from the joined errors.
		for _, je := range e.Unwrap() {
			return GetCode(je)
		}
	}

	// We really go out of our way to get an error code.
	if c := GetCode(errors.Unwrap(err)); c != "" {
		return c
	}
	return ECGeneric
}
