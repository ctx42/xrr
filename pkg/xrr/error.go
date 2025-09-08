package xrr

import (
	"encoding/json"
	"errors"
	"maps"
)

// WithCode is an option for [New] and [Wrap] setting the error code.
func WithCode(code string) func(*Error) {
	return func(e *Error) { e.code = code }
}

// Error represents an error with an error code and structured metadata.
type Error struct {
	error                // Wrapped error.
	code  string         // Error code.
	meta  map[string]any // Structured metadata.
}

// New creates a new [Error] instance with the specified message and error code.
// If the [WithCode] option is on the list of options, it will override the
// code argument.
func New(msg, code string, opts ...func(*Error)) error {
	err := &Error{
		error: errors.New(msg),
		code:  code,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

// Wrap wraps an error in an [Error] instance, applying the given options.
//
// It returns nil if the input error is nil or no options were provided. The
// returned error retains the same error code as the input error, obtained via
// [GetCode] function. To override the error code, use the [WithCode] option.
func Wrap(err error, opts ...func(*Error)) error {
	if err == nil {
		return nil
	}
	if len(opts) == 0 {
		return err
	}
	e := &Error{error: err, code: GetCode(err)}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// ErrorCode returns error code.
func (e *Error) ErrorCode() string { return e.code }

// MetaAll returns a clone of the error's metadata.
func (e *Error) MetaAll() map[string]any { return maps.Clone(e.meta) }

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.error
}

func (e *Error) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"error": e.Error(),
		"code":  e.code,
	}
	if len(e.meta) > 0 {
		m["meta"] = e.meta
	}
	return json.Marshal(m)
}
