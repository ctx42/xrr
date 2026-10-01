// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import "encoding/json"

// Compile time checks.
var (
	_ error          = Masked{}
	_ Coder          = Masked{}
	_ json.Marshaler = Masked{}
)

// Masked is an error that exposes the lead to callers while keeping the cause
// accessible in the Go error chain. `Error` and `ErrorCode` delegate to the
// lead; `MarshalJSON` renders only the lead. Both errors are reachable via
// [errors.Is], [errors.As], and the walk-based helpers ([GetMeta], etc.).
type Masked struct {
	cause error
	lead  error
}

// Mask creates a [Masked] from cause and lead. Returns nil if either cause or
// lead is nil.
func Mask(cause, lead error) error {
	if cause == nil || lead == nil {
		return nil
	}
	return Masked{cause: cause, lead: lead}
}

func (msk Masked) Error() string     { return msk.lead.Error() }
func (msk Masked) ErrorCode() string { return GetCode(msk.lead) }

// Cause returns the masked cause error, kept out of messages and JSON.
func (msk Masked) Cause() error { return msk.cause }

// Lead returns the error exposed to callers.
func (msk Masked) Lead() error { return msk.lead }

// MarshalJSON returns the JSON form of the lead error (using [errorToJSON]).
// The cause is intentionally not included in the output.
func (msk Masked) MarshalJSON() ([]byte, error) { return errorToJSON(msk.lead) }

func (msk Masked) Unwrap() []error {
	return []error{msk.cause, msk.lead}
}
