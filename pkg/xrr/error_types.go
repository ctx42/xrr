// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

// EDXrr is the marker type for the package's error domain.
type EDXrr struct{}

// Error constructor functions for the xrr package [edXrr] domain.
var (
	newError       = ErrorFunc[EDXrr]()
	newfError      = ErrorfFunc[EDXrr]()
	newFieldsError = FieldsFunc[EDXrr]()
)

// Error represents an error in the xrr package error domain.
type Error = GenericError[EDXrr]

// New creates a new [Error] with the given message. The args may contain an
// optional string error code and any number of [Option] values in any order.
//
// Examples:
//
//	xrr.New("message")
//	xrr.New("message", "ECode")
//	xrr.New("message", "ECode", xrr.Option())
//
// When [WithCause] is provided:
//   - If msg is empty, Error() returns the cause's message directly.
//   - If msg is non-empty, Error() returns "msg: cause message".
//   - If no code is given and [WithCode] is not provided, the cause's code is
//     inherited via [GetCode]. Pass a code string or [WithCode] to override it.
//
// For wrapping without a new message, prefer [Wrap], which makes the intent
// clearer.
func New(msg string, args ...any) error {
	return newError(msg, args...)
}

// Newf creates a new [Error] using a format string. It is the format-style
// counterpart of [New]: non-[Option] args are passed to the format string,
// while [Option] values are applied to the error. Unlike [New], a bare string
// argument is treated as a format argument, not an error code — pass [WithCode]
// to set the code.
//
// When the format string contains %w, the error is created via [fmt.Errorf]
// and stored as the cause; [GenericError.Error] delegates to it. Without %w,
// the message is set to fmt.Sprintf(format, args...).
//
// Examples:
//
//	xrr.Newf("user %d not found", userID)
//	xrr.Newf("user %d not found", userID, xrr.WithCode("ECode"))
//	xrr.Newf("connect failed: %w", err)
//	xrr.Newf("connect failed: %w", err, xrr.WithCode("ECode"))
func Newf(format string, args ...any) error {
	return newfError(format, args...)
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
