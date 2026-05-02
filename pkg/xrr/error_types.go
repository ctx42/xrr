// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

// EDXrr is the marker type for the package's error domain.
type EDXrr struct{}

// Error constructor functions for the xrr package [edXrr] domain.
var (
	newError       = ErrorFunc[EDXrr]()
	newFieldsError = FieldsFunc[EDXrr]()
)

// Error represents an error in the xrr package error domain.
type Error = GenericError[EDXrr]

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

// FieldErrors represents a field error in the xrr error domain.
type FieldErrors = GenericFields[EDXrr]

// NewFieldError returns a new [FieldErrors] containing the given field and
// error. Returns nil when the error is nil.
func NewFieldError(field string, err error) *FieldErrors {
	return newFieldsError(field, err)
}

// NewFieldErrors creates a new [FieldErrors] from the given map.
// The map is stored directly without copying.
func NewFieldErrors(fields map[string]error) *FieldErrors {
	return NewFields[EDXrr](fields)
}

// Wrap annotates err with a code and optional metadata, without adding a new
// message. The returned error's Error() is identical to err.Error(). It is
// [WrapUsing] bound to the default [EDXrr] domain.
//
// Returns nil if err is nil. The code defaults to the code of err (via
// [GetCode]); pass [WithCode] to override it. The original err is preserved in
// the error chain for [errors.Is] and [errors.As].
//
// To annotate with a new message as well, use [New] with [WithCause].
func Wrap(err error, opts ...Option) error {
	return WrapUsing[EDXrr](err, opts...)
}
