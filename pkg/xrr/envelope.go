// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
)

// Envelope provides facilities to create JSON envelope for errors.
//
// Cause is regular error and leading is not provided:
//
//	{
//	  "error": "cause",
//	  "code": "ECCause"
//	}
//
// Cause and leading are regular errors:
//
//	{
//	  "error": "main error",
//	  "code": "ECMain",
//	  "errors": [
//	    {"code": "ECCause", "error": "cause"},
//	  ]
//	}
//
// Cause is instance of [Fields] and leading error is provided:
//
//	{
//	  "error": "main",
//	  "code": "ECMain",
//	  "fields": {
//	    "field": {"code": "ECCause", "error": "cause"},
//	  }
//	}
//
// Cause is instance of [Fields] and leading error is not provided:
//
//	{
//	  "error": "fields error",
//	  "code": "ECFields",
//	  "fields": {
//	    "field": {"code": "ECCause", "error": "cause"},
//	  }
//	}
//
// Cause is join errors and leading error is provided:
//
//	{
//	  "error": "main error",
//	  "code": "ECMain",
//	  "errors": [
//	    "{"code":"ECE0","error":"msg 0"},
//	    "{"code":"ECE1","error":"msg 1"}
//	  ]
//	}
//
// Cause is join errors and leading error is not provided:
//
//	{
//	  "error": "msg 0",
//	  "code": "ECE0",
//	  "errors": [
//	    "{"code":"ECF1","error":"msg 1"}
//	  ]
//	}
type Envelope struct {
	cause error
	lead  error
}

// Enclose creates a new instance of [Envelope] from cause and optional leading
// error. Returns nil if the cause is nil. When more than one leading error is
// provided, only the first one is used.
func Enclose(cause error, lead ...error) error {
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

// Error returns the cause error message.
func (e Envelope) Error() string { return e.cause.Error() }

// ErrCode returns the cause error error code.
func (e Envelope) ErrCode() string { return GetCode(e.cause) }

// Unwrap returns the cause of the error.
func (e Envelope) Unwrap() error { return e.cause }

// Lead returns the leading error.
func (e Envelope) Lead() error { return e.lead }

// Is returns true if err is the same error as envelop or cause.
func (e Envelope) Is(target error) bool {
	return errors.Is(e.lead, target) || errors.Is(e.cause, target)
}

func (e Envelope) MarshalJSON() ([]byte, error) {
	var ef Fields
	if errors.As(e.cause, &ef) {
		if e.lead == nil {
			e.lead = ErrFields
		}
		return encloseFieldsError(e.lead, ef)
	}

	if IsJoined(e.cause) {
		ers := Split(e.cause)
		if e.lead == nil && len(ers) > 0 {
			e.lead = ers[0]
			ers = ers[1:]
		}
		return encloseMultiError(e.lead, ers...)
	}

	if e.lead != nil {
		return encloseMultiError(e.lead, e.cause)
	}
	return encloseMultiError(e.cause)
}

// encloseFieldsError returns [Fields] error enclosed in an error envelope with
// given leading error.
func encloseFieldsError(lead error, ef Fields) ([]byte, error) {
	ret := map[string]any{
		"error":  lead.Error(),
		"code":   GetCode(lead),
		"fields": ef,
	}
	if meta := GetMeta(lead); len(meta) > 0 {
		ret["meta"] = meta
	}
	return json.Marshal(ret)
}

// encloseMultiError returns multiple errors enclosed in an error envelope with
// given leading error.
func encloseMultiError(lead error, ers ...error) ([]byte, error) {
	ret := map[string]any{
		"error": lead.Error(),
		"code":  GetCode(lead),
	}
	if meta := GetMeta(lead); len(meta) > 0 {
		ret["meta"] = meta
	}
	if len(ers) > 0 {
		es := make([]json.RawMessage, len(ers))
		for i, e := range ers {
			entry, err := marshalError(e)
			if err != nil {
				return nil, err
			}
			es[i] = entry
		}
		ret["errors"] = es
	}
	return json.Marshal(ret)
}

// marshalError marshals error to JSON with check if the resulting JSON message
// is an empty object "{}", which means error did not have a [json.Marshaler]
// interface implemented, in which case the provided error is wrapped in the
// [Error] instance and marshaled again.
func marshalError(e error) ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	if len(data) == 2 {
		return json.Marshal(Wrap(e))
	}
	return data, nil
}
