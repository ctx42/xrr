// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"time"
	"unsafe"
)

// joined is an interface for an error that was created by [errors.Join].
type joined interface{ Unwrap() []error }

// Split splits joined errors into a slice of errors. It returns a single-element
// slice if the error does not implement `Unwrap() []error`. Returns nil if the
// input error is nil. Joined errors are typically created via [errors.Join].
func Split(err error) []error {
	if err == nil {
		return nil
	}
	if es, ok := err.(joined); ok {
		return es.Unwrap()
	}
	return []error{err}
}

// Join joins a slice of errors into a single error. Returns nil for an empty
// slice, or the single error if the slice has exactly one element. Otherwise
// delegates to [errors.Join].
func Join(ers ...error) error {
	ers = join(ers...)
	switch len(ers) {
	case 0:
		return nil
	case 1:
		return ers[0]
	default:
		return errors.Join(ers...)
	}
}

// join removes nil errors from the slice and returns the compacted result.
// It is the internal implementation backing the exported [Join] function.
//
// Empty input or all-nil input returns nil.
// The returned slice shares the backing array with the input (for efficiency),
// but is truncated to the number of non-nil errors.
func join(ers ...error) []error {
	if len(ers) == 0 {
		return nil
	}

	var j int
	for i := 0; i < len(ers); i++ {
		if err := ers[i]; err != nil {
			ers[j] = err
			j++
			continue
		}
	}
	ers = ers[:j]

	if len(ers) == 0 {
		return nil
	}
	return ers
}

// IsJoined returns true if the provided error is not nil and implements
// `Unwrap() []error` interface. Returns false if the error is nil.
func IsJoined(err error) bool {
	_, ok := err.(joined)
	return ok
}

// DefaultCode returns the first non-empty code from the slice of codes.
func DefaultCode(otherwise string, codes ...string) string {
	for _, code := range codes {
		if code != "" {
			return code
		}
	}
	return otherwise
}

// IsDomain returns true if err is a [GenericError] or [GenericFields] of
// domain T.
func IsDomain[T Domain](err error) bool {
	if _, ok := err.(*GenericError[T]); ok {
		return true
	}
	if _, ok := err.(*GenericFields[T]); ok {
		return true
	}
	return false
}

// isNil reports whether v is nil, including the case of a typed nil
// (e.g. (*T)(nil), (error)(nil), etc.).
//
// This is needed because a nil interface holding a typed nil is not equal
// to a bare nil interface in Go. It is used as a defensive guard in the
// common pattern:
//
//	if x == nil || isNil(x) { ... }
//
// Callers include [GetCode], [walkErrors], [walkErrorsReverse], and
// [WrapUsing].
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Func,
		reflect.Slice, reflect.Map, reflect.Chan:
		return rv.IsNil()
	default:
		return false
	}
}

// joinKey joins a prefix and key with a dot when the prefix is non-empty.
// It is used internally to build dotted paths when flattening nested
// [Fielder] structures (see [flattenFieldErrors] and [GenericFields]).
func joinKey(pref, key string) string {
	if pref != "" {
		if key == "" {
			return pref
		}
		return pref + "." + key
	}
	return key
}

// isValidMetaValue reports whether v has a type that is allowed to be stored
// as metadata. This is the runtime check that matches the [MetaType]
// constraint.
//
// Supported types: bool, string, int, int64, float64, [time.Time],
// [time.Duration]. Unsupported values are silently skipped by [WithMeta],
// [Metadata.MetaSetAll], and [Metadata.MetaSetFrom].
func isValidMetaValue(v any) bool {
	switch v.(type) {
	case bool, string, int, int64, float64, time.Time, time.Duration:
		return true
	default:
		return false
	}
}

// sortFieldErrors converts a map of field errors to two parallel slices
// (field names and errors), with the names sorted in ascending order.
// It is used by [formatFieldErrors] and the [walkErrors] family.
func sortFieldErrors(ers map[string]error) ([]string, []error) {
	var errs []error
	var fields []string
	for field := range ers {
		fields = append(fields, field)
	}
	slices.Sort(fields)
	for _, field := range fields {
		errs = append(errs, ers[field])
	}
	return fields, errs
}

// errorMessage returns the message portion of an error for use in
// [GenericError.Error] and field error formatting.
//
// For joined errors (those implementing `Unwrap() []error`):
//   - A single error returns its message directly.
//   - Multiple errors have their messages concatenated with "; ".
//   - All other errors return err.Error() directly.
//
// The unsafe.String optimization avoids an extra allocation when building
// the concatenated message for joined errors.
func errorMessage(err error) string {
	if jes, ok := err.(joined); ok {
		es := jes.Unwrap()
		if len(es) == 1 {
			return es[0].Error()
		}

		// Build "msg1; msg2; msg3..." efficiently without extra allocations.
		b := []byte(es[0].Error())
		for _, e := range es[1:] {
			b = append(b, ';', ' ')
			b = append(b, e.Error()...)
		}
		// b is guaranteed to have at least two bytes here.
		return unsafe.String(&b[0], len(b)) // nolint: gosec
	}
	return err.Error()
}

// errorToJSON returns the JSON representation of an error according to the
// library's conventions.
//
// It first tries the error's own [json.Marshaler] implementation. If that
// produces an empty JSON object, it falls back to [errorToMap] to build the
// standard {"error", "code", "meta"} form. A marshal error is returned
// directly.
//
// This is the function used by [Masked.MarshalJSON].
func errorToJSON(e error) ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	if len(data) == 2 {
		return json.Marshal(errorToMap(e))
	}
	return data, nil
}

// errorToMap returns the canonical map form used by the library for JSON
// serialization of errors: {"error", "code", "meta"}.
//
// Returns nil if err is nil. This is the fallback used by [errorToJSON]
// when an error does not have a useful [json.Marshaler] implementation.
// It is also used directly by [envelopeFieldsToJSON] and [envelopeErrorsToJSON].
func errorToMap(err error) map[string]any {
	if err == nil {
		return nil
	}
	m := map[string]any{
		"error": err.Error(),
		"code":  GetCode(err),
	}
	if meta := GetMeta(err); len(meta) > 0 {
		m["meta"] = meta
	}
	return m
}

// newArgs extracts an optional error code and zero or more [Option] values
// from a mixed argument list.
//
// Rules:
//   - The last string argument is treated as the error code (earlier strings
//     are ignored).
//   - All [Option] values are collected in order.
//   - Any other types are ignored.
//
// This enables the flexible argument style of [ErrorFunc], [New], etc.
func newArgs(ags ...any) (string, []Option) {
	var code string
	var opts []Option
	for _, arg := range ags {
		switch a := arg.(type) {
		case string:
			code = a // last one wins
		case Option:
			opts = append(opts, a)
		default:
			// Intentionally ignored (see doc comment).
		}
	}
	return code, opts
}

// newfArgs separates format-style constructor arguments into three groups.
//
// It returns:
//   - wrapsErr: true if the format contains %w (the first non-Option error
//     arg will be turned into a cause via [WithCause]).
//   - args: non-[Option] values (including errors) for fmt.Sprintf / fmt.Errorf.
//   - opts: all [Option] values, collected separately.
//
// Used by [ErrorfFunc] and [Newf].
func newfArgs(format string, ags ...any) (bool, []any, []Option) {
	var args []any
	var opts []Option
	wrapsErr := strings.Index(format, "%w") >= 0
	for _, arg := range ags {
		switch a := arg.(type) {
		case Option:
			opts = append(opts, a)
		default:
			// Both errors and everything else are treated as format arguments.
			// When wrapsErr is true, the first error will later be turned into
			// a cause by ErrorfFunc.
			args = append(args, a)
		}
	}
	return wrapsErr, args, opts
}
