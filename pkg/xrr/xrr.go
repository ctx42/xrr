// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package xrr provides errors supporting error codes and metadata.
package xrr

// ECGeneric represents generic error code used for non-nil errors which have
// no error code assigned.
const ECGeneric = "ECGeneric"

// Domain represents types that can be used to define error domains.
type Domain interface{ comparable }

// EDGeneric represents generic error domain used in the xrr package.
type EDGeneric string

// Coder is the interface that wraps the ErrorCode method.
type Coder interface {
	// ErrorCode returns the error code for the error. For errors without an
	// explicit code, it should return [ECGeneric]. Callers should prefer
	// [GetCode], which handles nil errors before invoking this method.
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

// Error constructor functions for the xrr package [EDGeneric] domain.
var (
	newError    = ErrorFactory[EDGeneric]()
	fieldsError = FieldsFactory[EDGeneric]()
)

// New creates a new [GenericError[EDGeneric]] error instance with the given
// message and error code. If the [WithCode] option is on the list of
// options, it will override the code argument.
func New(msg, code string, opts ...Option) error {
	return newError(msg, code, opts...)
}

// NewField returns a new instance of [Fields] with the given error. Returns
// nil when the error is nil.
func NewField(field string, err error) error {
	return fieldsError(field, err)
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
