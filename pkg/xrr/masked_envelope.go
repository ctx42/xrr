// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

// MaskedEnvelope is an error that exposes the lead to callers while keeping
// the cause accessible in the Go error chain. [Error] and [ErrorCode] delegate
// to the lead; [MarshalJSON] renders only the lead. Both errors are reachable
// via [errors.Is], [errors.As], and the walk-based helpers ([GetMeta], etc.).
type MaskedEnvelope struct {
	cause error
	lead  error
}

// EncloseMasked creates a [MaskedEnvelope] from cause and lead. Returns nil if
// either cause or lead is nil.
func EncloseMasked(cause, lead error) error {
	if cause == nil || lead == nil {
		return nil
	}
	return MaskedEnvelope{cause: cause, lead: lead}
}

func (e MaskedEnvelope) Error() string                { return e.lead.Error() }
func (e MaskedEnvelope) ErrorCode() string            { return GetCode(e.lead) }
func (e MaskedEnvelope) Lead() error                  { return e.lead }
func (e MaskedEnvelope) Unwrap() []error              { return []error{e.cause, e.lead} }
func (e MaskedEnvelope) MarshalJSON() ([]byte, error) { return marshalError(e.lead) }
