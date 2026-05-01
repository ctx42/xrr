// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package xrr provides errors supporting error codes and metadata.
package xrr

// ECGeneric represents generic error code used for non-nil errors which have
// no error code assigned.
const ECGeneric = "ECGeneric"

// Domain represents types that can be used to define error domains.
type Domain interface{ comparable }

// edXrr is the marker type for the package's error domain.
type edXrr struct{}

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

// Error constructor functions for the xrr package [edXrr] domain.
var (
	newError       = ErrorFunc[edXrr]()
	newFieldsError = FieldsFunc[edXrr]()
)

// Error represents an error in the verax package error domain.
type Error = GenericError[edXrr]

// New creates a new [Error] with the given message and error code.
//
// When [WithCause] is provided:
//   - If msg is empty, Error() returns the cause's message directly.
//   - If msg is non-empty, Error() returns "msg: cause message".
//   - If code is empty and [WithCode] is not provided, the cause's code is
//     inherited via [GetCode]. Pass a non-empty code argument or [WithCode]
//     to override it.
//
// For wrapping without a new message, prefer [Wrap] which makes the intent
// clearer.
func New(msg, code string, opts ...Option) error {
	return newError(msg, code, opts...)
}

// FieldErrors represents a field error in the verax error domain.
type FieldErrors = GenericFields[edXrr]

// NewFieldError returns a new [FieldErrors] containing the given field and
// error. Returns nil when the error is nil.
func NewFieldError(field string, err error) *FieldErrors {
	return newFieldsError(field, err)
}

// NewFieldErrors creates a new [FieldErrors] from the given map.
// The map is stored directly without copying.
func NewFieldErrors(fields map[string]error) *FieldErrors {
	return NewFields[edXrr](fields)
}

// Wrap wraps an error in a [GenericError[T]] instance, applying the given
// options. The wrapped error is accessible via [errors.Unwrap] and participates
// in [errors.Is] / [errors.As] chain traversal.
//
// Returns nil if err is nil. The returned error inherits the code of err via
// [GetCode]; use [WithCode] to override it. [Error] returns the cause's message
// directly (no new message is associated with the wrapper).
//
// When both a new message and a cause are needed in the same error, use
// [New] with [WithCause] instead.
func Wrap[T Domain](err error, opts ...Option) error {
	if err == nil || isNil(err) {
		return nil
	}
	ops := Options{code: GetCode(err)}.Set(opts...)
	return &GenericError[T]{
		msg:  "",
		code: ops.code,
		meta: ops.meta,
		err:  err,
	}
}

// SetCode assigns code to err by wrapping it with [Wrap]. Returns nil if err
// is nil. Returns err unchanged if code is empty or err already carries the
// given code.
func SetCode[T Domain](err error, code string) error {
	if code == "" {
		return err
	}
	if have := GetCode(err); have == code {
		return err
	}
	return Wrap[T](err, WithCode(code))
}
