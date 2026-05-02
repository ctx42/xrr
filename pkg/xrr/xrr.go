// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package xrr provides errors supporting error codes and metadata.
package xrr

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

// WrapUsing annotates err with a code and optional metadata in domain T,
// without adding a new message. The returned error's Error() is identical to
// err.Error().
//
// Returns nil if err is nil. The code defaults to the code of err (via
// [GetCode]); pass [WithCode] to override it. The original err is preserved in
// the error chain for [errors.Is] and [errors.As].
//
// To annotate with a new message as well, obtain a constructor with [ErrorFunc]
// and pass [WithCause].
func WrapUsing[T Domain](err error, opts ...Option) error {
	if err == nil || isNil(err) {
		return nil
	}
	ops := Options{code: GetCode(err)}.Set(opts...)
	return &GenericError[T]{
		code: ops.code,
		meta: ops.meta,
		err:  err,
	}
}

// SetCode assigns code to err by wrapping it with [WrapUsing]. Returns nil if
// err is nil. Returns err unchanged if code is empty or err already carries
// the given code.
func SetCode[T Domain](err error, code string) error {
	if code == "" {
		return err
	}
	if have := GetCode(err); have == code {
		return err
	}
	return WrapUsing[T](err, WithCode(code))
}
