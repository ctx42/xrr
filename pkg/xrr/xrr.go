// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package xrr extends the standard library's error handling with stable
// error codes, structured metadata, and domain-specific error types using
// Go generics.
//
// It is fully compatible with [errors.Is], [errors.As], [errors.Join], and
// [fmt.Errorf] wrapping. The library never replaces the error interface.
//
// Core concepts:
//
//   - Error codes via the [Coder] interface (stable identifiers for
//     monitoring, API responses, etc.).
//   - Structured metadata via [Metadater] and the [Meta] builder.
//   - Domain-specific errors using generics (e.g. PaymentError) for
//     type-safe error handling across subsystems.
//   - Field validation errors via [Fielder] and [NewFieldErrors].
//   - Rich inspection with [GetCode], [GetMeta], [IsCode], etc.
//   - First-class support for JSON marshaling and slog integration.
//
// See the README for usage examples.
package xrr

import (
	"fmt"
)

// ECGeneric represents generic error code used for non-nil errors which have
// no error code assigned.
const ECGeneric = "ECGeneric"

// Domain represents types that can be used to define error domains.
type Domain interface{ comparable }

// Coder is the interface that wraps the ErrorCode method.
type Coder interface {
	// ErrorCode returns the error code for the error. For errors without an
	// explicit code, it should return [ECGeneric].
	ErrorCode() string
}

// Fielder is the interface that wraps the ErrorFields method.
//
// ErrorFields returns errors for field names. It is used for validation errors.
type Fielder interface {
	ErrorFields() map[string]error
}

// Metadater is an interface providing access to error metadata.
type Metadater interface {
	// MetaAll returns a copy of the metadata held directly by this error.
	//
	// It does not include metadata from wrapped errors. To retrieve metadata
	// recursively, use [GetMeta] instead. Returns nil if no metadata is
	// present.
	MetaAll() map[string]any
}

// WrapUsing annotates err with a code and/or metadata in domain T, without
// changing its message (unless [WithCause] is used — see below).
//
// The original err is preserved as the primary value in the error chain for
// [errors.Is] and [errors.As].
//
// When [WithCause] is supplied, the resulting error's Error() string becomes
// "original; cause". The original err remains the direct target of
// [errors.Unwrap]; both the original and the cause remain discoverable via
// [errors.Is] and [errors.As].
//
// Returns nil if err is nil. The code defaults to the code of err (via
// [GetCode]); pass [WithCode] to override it.
//
// To create a fresh error with a message + cause, use [ErrorFunc] + [WithCause]
// instead.
//
// Example (basic annotation):
//
//	err := someOperation()
//	if err != nil {
//	    return xrr.WrapUsing[MyDomain](err, xrr.WithCode("E_FOO"))
//	}
//
// Example with additional cause (note that both the original error and the
// cause remain discoverable via [errors.Is]):
//
//	cause := errors.New("db connection failed")
//	highLevelErr := someHighLevelOperation()
//
//	wrapped := xrr.WrapUsing[MyDomain](highLevelErr, xrr.WithCause(cause))
//
//	fmt.Println(errors.Is(wrapped, cause))       // true
//	fmt.Println(errors.Is(wrapped, highLevelErr)) // true
func WrapUsing[T Domain](err error, opts ...Option) error {
	if err == nil || isNil(err) {
		return nil
	}

	ops := Options{code: GetCode(err)}.Set(opts...)

	if ops.err != nil {
		err = withAdditionalCause(err, ops.err)
	}

	return &GenericError[T]{
		code: ops.code,
		meta: ops.meta,
		err:  err,
	}
}

// withAdditionalCause attaches an additional cause while keeping the
// original error as the direct target of [errors.Unwrap].
//
// The returned error's [errors.Unwrap] yields the original err (not the
// cause). Both the original and the cause remain discoverable via
// [errors.Is] and [errors.As]. The Error() string becomes "original; cause".
//
// This internal helper exists only to support the special [WithCause]
// behavior inside [WrapUsing].
func withAdditionalCause(err, cause error) error {
	if cause == nil {
		return err
	}
	if err == nil {
		return cause
	}
	return fmt.Errorf("%w; %w", err, cause)
}

// SetCode assigns a code to err by wrapping it with [WrapUsing].
// It is a convenience over calling WrapUsing directly with [WithCode].
//
// Returns nil if err is nil. Returns err unchanged if code is empty or
// err already carries the given code.
func SetCode[T Domain](err error, code string) error {
	if code == "" {
		return err
	}
	if have := GetCode(err); have == code {
		return err
	}
	return WrapUsing[T](err, WithCode(code))
}
