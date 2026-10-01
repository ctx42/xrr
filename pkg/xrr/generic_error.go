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
// The message is built with [fmt.Errorf]. When it wraps an error through a %w
// verb (including indexed forms such as %[1]w), the result becomes the cause
// (preserving the error chain for [errors.Is] / [errors.As]) and the message
// is left empty so [GenericError.Error] delegates to it. Otherwise the
// formatted text becomes the message.
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
		args, opts := newfArgs(args...)

		var msg string
		if err := fmt.Errorf(format, args...); wrapsError(err) {
			opts = append(opts, WithCause(err))
		} else {
			msg = err.Error()
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
func (gen *GenericError[T]) Error() string {
	if gen.err != nil {
		em := errorMessage(gen.err)
		if gen.msg != "" {
			return gen.msg + ": " + em
		}
		return em
	}
	return gen.msg
}

// ErrorCode returns error code. Returns [ECGeneric] when no code is set.
func (gen *GenericError[T]) ErrorCode() string {
	if gen.code == "" {
		return ECGeneric
	}
	return gen.code
}

// MetaAll returns a clone of the error's metadata.
func (gen *GenericError[T]) MetaAll() map[string]any { return maps.Clone(gen.meta) }

// Unwrap returns the wrapped error.
func (gen *GenericError[T]) Unwrap() error {
	if gen == nil {
		return nil
	}
	return gen.err
}

func (gen *GenericError[T]) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"error": gen.Error(),
		"code":  gen.ErrorCode(),
	}
	if meta := GetMeta(gen); len(meta) > 0 {
		m["meta"] = meta
	}
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals the JSON representation of a [GenericError].
//
// The minimal valid form is {"error": "message"}; in that case the code
// defaults to [ECGeneric].
//
// The JSON form carries no cause, so any wrapped error the receiver held
// before is cleared.
//
// Note: Numeric metadata values are unmarshaled as float64.
func (gen *GenericError[T]) UnmarshalJSON(data []byte) error {
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

	gen.msg = msg
	gen.code = code
	gen.meta = meta
	gen.err = nil
	return nil
}

// Format implements [fmt.Formatter] for [GenericError].
func (gen *GenericError[T]) Format(state fmt.State, verb rune) {
	Format(gen.Error(), gen.ErrorCode(), state, verb)
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
