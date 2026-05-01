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

// ErrorFunc returns a function for creating domain-specific errors.
func ErrorFunc[T Domain]() func(msg, code string, opts ...Option) *GenericError[T] {
	return func(msg, code string, opts ...Option) *GenericError[T] {
		ops := Options{code: code}.Set(opts...)
		return &GenericError[T]{
			msg:  msg,
			code: ops.code,
			meta: ops.meta,
			err:  ops.err,
		}
	}
}

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

// UnmarshalJSON unmarshals JSON representation of the [GenericError].
//
// The minimal valid JSON representation for a [GenericError] is
//
//	{"error": "message"}
//
// and in this case, the error code is set to [ECGeneric].
//
// Notes:
//   - Numeric values will be unmarshalled as float64.
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

// Format is a custom formatter for [GenericError] instances.
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
