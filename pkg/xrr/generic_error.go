// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Compile time checks.
var (
	_ error            = (*GenericError[EDXrr])(nil)
	_ Coder            = (*GenericError[EDXrr])(nil)
	_ Metadater        = (*GenericError[EDXrr])(nil)
	_ json.Marshaler   = (*GenericError[EDXrr])(nil)
	_ json.Unmarshaler = (*GenericError[EDXrr])(nil)
)

// GenericError represents a generic type for creating domain-specific errors.
type GenericError[T Domain] struct {
	msg  string         // Error message.
	code string         // Error code.
	meta map[string]any // Structured metadata.
	err  error          // Wrapped error.
}

// ErrorFunc returns a constructor for domain-T errors. The returned function
// accepts a message string followed by an optional error code string and any
// number of [Option] values in any order.
//
// Examples:
//
//	newErr := xrr.ErrorFunc[MyDomain]()
//	newErr("message")
//	newErr("message", "ECode")
//	newErr("message", "ECode", xrr.Option())
func ErrorFunc[T Domain]() func(msg string, args ...any) *GenericError[T] {
	return func(msg string, args ...any) *GenericError[T] {
		code, opts := newArgs(args...)
		ops := Options{code: code}.Set(opts...)
		return &GenericError[T]{
			msg:  msg,
			code: ops.code,
			meta: ops.meta,
			err:  ops.err,
		}
	}
}

// ErrorfFunc returns a format-style constructor for domain-T errors. The
// returned function behaves like [fmt.Sprintf]: non-[Option] args are passed
// as format arguments, and [Option] values are applied to the error.
//
// When the format string contains %w, the constructor uses [fmt.Errorf] to
// build the cause (preserving the error chain for [errors.Is] / [errors.As])
// and leaves the message empty so [GenericError.Error] delegates to the cause.
// Without %w, the message is set to fmt.Sprintf(format, args...).
//
// The error code must be passed via [WithCode]; unlike [ErrorFunc], a bare
// string argument is treated as a format argument, not a code.
//
// Examples:
//
//	newErrf := xrr.ErrorfFunc[MyDomain]()
//	newErrf("user %d not found", userID)
//	newErrf("user %d not found", userID, xrr.WithCode("ECode"))
//	newErrf("connect failed: %w", err)
//	newErrf("connect failed: %w", err, xrr.WithCode("ECode"))
func ErrorfFunc[T Domain]() func(format string, args ...any) *GenericError[T] {
	return func(format string, args ...any) *GenericError[T] {
		wraps, args, opts := newfArgs(format, args...)

		var msg string
		if wraps {
			err := fmt.Errorf(format, args...)
			opts = append(opts, WithCause(err))
		} else {
			msg = fmt.Sprintf(format, args...)
		}

		ops := Options{}.Set(opts...)
		return &GenericError[T]{
			msg:  msg,
			code: ops.code,
			meta: ops.meta,
			err:  ops.err,
		}
	}
}

// Error returns the error message. For errors with a cause or joined errors
// it delegates to [errorMessage] (which produces "; "-separated messages for
// joined errors).
func (e *GenericError[T]) Error() string {
	if e.err != nil {
		em := errorMessage(e.err)
		if e.msg != "" {
			return e.msg + ": " + em
		}
		return em
	}
	return e.msg
}

// ErrorCode returns error code. Returns [ECGeneric] when no code is set.
func (e *GenericError[T]) ErrorCode() string {
	if e.code == "" {
		return ECGeneric
	}
	return e.code
}

// MetaAll returns a clone of the error's metadata.
func (e *GenericError[T]) MetaAll() map[string]any { return maps.Clone(e.meta) }

// Unwrap returns the wrapped error.
func (e *GenericError[T]) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *GenericError[T]) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"error": e.Error(),
		"code":  e.ErrorCode(),
	}
	if meta := GetMeta(e); len(meta) > 0 {
		m["meta"] = meta
	}
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals the JSON representation of a [GenericError].
//
// The minimal valid form is {"error": "message"}; in that case the code
// defaults to [ECGeneric].
//
// Note: Numeric metadata values are unmarshaled as float64.
func (e *GenericError[T]) UnmarshalJSON(data []byte) error {
	m := make(map[string]any, 3)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	msgI := m["error"]
	msg, _ := msgI.(string)
	if msg == "" {
		return ErrInvJSONError
	}

	codeI := m["code"]
	code, _ := codeI.(string)
	if code == "" {
		code = ECGeneric
	}

	metaI := m["meta"]
	var meta map[string]any
	if metaI != nil {
		meta, _ = metaI.(map[string]any)
	}

	e.msg = msg
	e.code = code
	e.meta = meta
	return nil
}

// Format implements [fmt.Formatter] for [GenericError].
func (e *GenericError[T]) Format(state fmt.State, verb rune) {
	Format(e.Error(), e.ErrorCode(), state, verb)
}

// Format is the implementation backing [GenericError.Format] (and used by
// field error formatting).
func Format(msg, code string, state fmt.State, verb rune) {
	switch verb {
	case 's', 'q':
		if verb == 'q' {
			msg = fmt.Sprintf("%q", msg)
		}
		_, _ = fmt.Fprint(state, msg)

	case 'v':
		_, _ = fmt.Fprint(state, msg)
		if state.Flag('+') {
			_, _ = fmt.Fprintf(state, " (%s)", code)
		}
	}
}
