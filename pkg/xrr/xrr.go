// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package xrr provides errors supporting error codes and metadata.
package xrr

// ECGeneric represents generic error code used for non-nil errors which have
// no error code assigned.
const ECGeneric = "ECGeneric"

// Domain represents types that can be used to define error domains.
type Domain interface{ ~string }

// EDGeneric represents generic error domain used in the xrr package.
type EDGeneric string

// Coder is the interface that wraps the ErrorCode method.
//
// ErrorCode returns error code.
//
// For nil errors it must return an empty string, but for non-nil errors
// without assigned code, it should return [ECGeneric] error code.
type Coder interface {
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
	// MetaAll returns a copy of the metadata from the [GenericError] instance.
	//
	// It does not include metadata from wrapped errors. To retrieve metadata,
	// recursively, use the [GetMeta] function instead. Returns nil if no
	// metadata is present.
	MetaAll() map[string]any
}

// Error constructor functions for the xrr package [EDGeneric] domain.
var (
	newError   = NewErrorFor[EDGeneric]()
	fieldError = NewFieldErrorFor[EDGeneric]()
)

// Error represents an error type in [EDGeneric] domain.
// type Error = GenericError[EDGeneric]

// New creates a new [GenericError[EDGeneric]] error instance with the given
// message and error code. If the [WithCode] option is on the list of
// options, it will override the code argument.
func New(msg, code string, opts ...Option) error {
	return newError(msg, code, opts...)
}

// Fields represents a collection of errors that are indexed by field names.
//
// The field (name) indexed errors are mostly useful for validation errors.
// type Fields = GenericFields[EDGeneric]

// FieldError returns a new instance of [Fields] with the given error. Returns
// nil when the error is nil.
func FieldError(field string, err error) error {
	return fieldError(field, err)
}

// Wrap wraps an error in a [GenericError] instance, applying the given options.
//
// It returns nil if the error is nil. The returned error retains the same
// error code as the input error, obtained via [GetCode] function. To override
// the error code, use the [WithCode] option.
func Wrap[T Domain](err error, opts ...Option) error {
	if err == nil {
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
