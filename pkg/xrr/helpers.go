package xrr

import (
	"reflect"
)

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
