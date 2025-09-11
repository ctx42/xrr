package xrr

import (
	"errors"
	"reflect"
)

// Split splits joined errors into a slice of. It will return the slice with a
// single error if the provided error does not implement the `Unwrap []error`
// interface. It will return nil if the error is nil.
func Split(err error) []error {
	if err == nil {
		return nil
	}
	var joinErr interface{ Unwrap() []error }
	if errors.As(err, &joinErr) {
		return joinErr.Unwrap()
	}
	return []error{err}
}

// DefaultCode returns the first non-empty code from the slice of codes.
func DefaultCode(otherwise string, codes ...string) string {
	for _, code := range codes {
		if code != "" {
			return code
		}
	}
	return otherwise
}

// isNil returns true if v is nil or v is nil interface.
func isNil(v any) bool {
	defer func() { _ = recover() }()
	return v == nil || reflect.ValueOf(v).IsNil()
}
