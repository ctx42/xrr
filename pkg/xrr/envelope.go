// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
)

// Envelope provides a JSON envelope for errors.
//
// An [Envelope] has two fields: `cause` (the error encountered during
// execution) and `lead` (the error to surface as the top-level error in the
// JSON output). See the examples below.
//
// - the `lead` error is not provided:
//
//	{
//	  "error": "cause",
//	  "code": "ECCause"
//	}
//
// - the `cause` and `lead` are provided:
//
//	{
//	  "error": "lead",
//	  "code": "ECLead",
//	  "errors": [
//	    {"code": "ECCause", "error": "cause"},
//	  ]
//	}
//
// - the `cause` is an instance of [Fielder] with `lead` error provided:
//
//	{
//	  "error": "lead",
//	  "code": "ECLead",
//	  "fields": {
//	    "field": {"code": "ECCause", "error": "cause"},
//	  }
//	}
//
// - the `cause` is a [Fielder] and `lead` error is not provided:
//
//	{
//	  "error": "fields error",
//	  "code": "ECFields",
//	  "fields": {
//	    "field": {"code": "ECCause", "error": "cause"},
//	  }
//	}
//
// - the `cause` is a joined error and `lead` error is provided:
//
//	{
//	  "error": "lead",
//	  "code": "ECLead",
//	  "errors": [
//	    {"code":"ECCause0","error":"cause 0"},
//	    {"code":"ECCause1","error":"cause 1"}
//	  ]
//	}
//
// - the `cause` is a joined error and `lead` error is not provided:
//
//	{
//	  "error": "cause 0",
//	  "code": "ECCause0",
//	  "errors": [
//	    {"code":"ECCause1","error":"cause 1"}
//	  ]
//	}
type Envelope struct {
	cause error
	lead  error
}

// Envelop creates a new instance of [Envelope] from cause and optional leading
// error. Returns nil if the cause is nil. When more than one leading error is
// provided, only the first one is used.
func Envelop(cause error, lead ...error) error {
	if cause == nil {
		return nil
	}

	// If cause is already an instance of Envelope,
	// use it and override the leading error if provided.
	//
	// nolint: errorlint
	//goland:noinspection ALL
	if e, ok := cause.(Envelope); ok {
		if len(lead) > 0 {
			e.lead = lead[0]
		}
		return e
	}

	enc := Envelope{cause: cause}
	if len(lead) > 0 {
		enc.lead = lead[0]
	}
	return enc
}

func (e Envelope) Error() string { return e.cause.Error() }

func (e Envelope) ErrorCode() string { return GetCode(e.cause) }

// Unwrap returns the cause of the error.
func (e Envelope) Unwrap() error { return e.cause }

// Lead returns the leading error.
func (e Envelope) Lead() error { return e.lead }

// Is returns true if target matches the lead or cause error.
func (e Envelope) Is(target error) bool {
	return errors.Is(e.lead, target) || errors.Is(e.cause, target)
}

// MarshalJSON returns the JSON representation of the envelope.
//
// The exact shape depends on the contents:
//   - If the cause is a [Fielder], it delegates to [envelopeFieldsToJSON].
//   - If the cause is a joined error, it delegates to [envelopeErrorsToJSON].
//   - Otherwise it uses the lead (if present) + cause via [envelopeErrorsToJSON].
func (e Envelope) MarshalJSON() ([]byte, error) {
	if ef, ok := e.cause.(Fielder); ok {
		if e.lead == nil {
			e.lead = ErrFields
		}
		return envelopeFieldsToJSON(e.lead, ef)
	}

	if IsJoined(e.cause) {
		ers := Split(e.cause)
		if e.lead == nil && len(ers) > 0 {
			e.lead = ers[0]
			ers = ers[1:]
		}
		return envelopeErrorsToJSON(e.lead, ers...)
	}

	if e.lead != nil {
		return envelopeErrorsToJSON(e.lead, e.cause)
	}
	return envelopeErrorsToJSON(e.cause)
}

// envelopeFieldsToJSON returns the JSON representation of an
// [Envelope] that contains field validation errors.
//
// It is called by [Envelope.MarshalJSON] when the cause implements [Fielder].
// The resulting object includes the lead error plus a "fields" key.
func envelopeFieldsToJSON(lead error, ef Fielder) ([]byte, error) {
	ret := errorToMap(lead)
	ret["fields"] = &GenericFields[EDXrr]{fields: ef.ErrorFields()}
	if meta := GetMeta(lead); len(meta) > 0 {
		ret["meta"] = meta
	}
	return json.Marshal(ret)
}

// envelopeErrorsToJSON returns the JSON representation of an
// [Envelope] that wraps multiple errors (or a single cause with a lead).
//
// It is called by [Envelope.MarshalJSON] for joined errors, for a lead
// error plus one or more causes, and for plain single-cause envelopes.
// The resulting object includes an "errors" array when there is more than
// one error.
func envelopeErrorsToJSON(lead error, ers ...error) ([]byte, error) {
	ret := errorToMap(lead)
	if len(ers) > 0 {
		es := make([]json.RawMessage, len(ers))
		for i, e := range ers {
			entry, err := errorToJSON(e)
			if err != nil {
				return nil, err
			}
			es[i] = entry
		}
		ret["errors"] = es
	}
	return json.Marshal(ret)
}
