package xrr

import (
	"errors"
)

// Error represents an error with an error code and structured metadata.
type Error struct {
	error                // Wrapped error.
	code  string         // Error code.
	meta  map[string]any // Structured metadata.
}

// New creates a new [Error] instance with the specified message and error code.
func New(msg, code string, opts ...func(*Error)) *Error {
	err := &Error{
		error: errors.New(msg),
		code:  code,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
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
