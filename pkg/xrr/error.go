package xrr

import (
	"errors"
)

// Error represents an error with an optional erro code and structured metadata.
type Error struct {
	error                // Wrapped error.
	code  string         // Error code.
	meta  map[string]any // Structured metadata.
}

// New creates a new [Error] instance with the specified error code and message.
//
// If no error code is provided, [ECGeneric] is used as the default.
// If multiple error codes are provided, only the first is used.
func New(msg string, code ...string) *Error {
	return &Error{error: errors.New(msg), code: DefaultCode(ECGeneric, code...)}
}

// ErrorCode returns error code.
func (e *Error) ErrorCode() string { return e.code }

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.error
}
