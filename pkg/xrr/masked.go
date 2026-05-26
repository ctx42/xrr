// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import "encoding/json"

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

func (e Masked) Error() string     { return e.lead.Error() }
func (e Masked) ErrorCode() string { return GetCode(e.lead) }
func (e Masked) Cause() error      { return e.cause }
func (e Masked) Lead() error       { return e.lead }

// MarshalJSON returns the JSON form of the lead error (using [errorToJSON]).
// The cause is intentionally not included in the output.
func (e Masked) MarshalJSON() ([]byte, error) { return errorToJSON(e.lead) }

func (e Masked) Unwrap() []error {
	return []error{e.cause, e.lead}
}

// Compile time checks.
var _ json.Marshaler = Masked{}
