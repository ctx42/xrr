// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
	"reflect"
	"sort"
	"time"
	"unsafe"
)

// joined is an interface for an error that was created by [errors.Join].
type joined interface{ Unwrap() []error }

// Split splits joined errors into a slice of. It will return the slice with a
// single error if the provided error does not implement the `Unwrap []error`
// interface. It will return nil if the error is nil.
func Split(err error) []error {
	if err == nil {
		return nil
	}
	if es, ok := err.(joined); ok {
		return es.Unwrap()
	}
	return []error{err}
}

// Join joins a slice of errors into a single error. It will return nil if the
// slice is empty. It will return the single error if the slice contains only
// one error. Otherwise, it will use [errors.Join] to join the errors.
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

// isNil returns true if v is nil or v is nil interface.
func isNil(v any) bool {
	defer func() { _ = recover() }()
	return v == nil || reflect.ValueOf(v).IsNil()
}

// prefix adds prefix to the key if the prefix is not empty.
func prefix(pref, key string) string {
	if pref != "" {
		if key == "" {
			return pref
		}
		return pref + "." + key
	}
	return key
}

// isTypeSupported returns true if the type of v is the supported metadata type.
func isTypeSupported(v any) bool {
	switch v.(type) {
	case bool, string, int, int64, float64, time.Time, time.Duration:
		return true
	default:
		return false
	}
}

// sortFields converts a map of errors to two slices: one for field names and
// one for errors. The returned slices maintain corresponding indexes, ensuring
// that each field name aligns with its associated error. Both slices are
// always of equal length, and the field names are sorted in ascending order.
func sortFields(ers map[string]error) ([]string, []error) {
	var errs []error
	var fields []string
	for field := range ers {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		errs = append(errs, ers[field])
	}
	return fields, errs
}

// errorMessage formats the error message for the given error. If the error
// implements the `Unwrap() []error` interface it concatenates the messages of
// all unwrapped errors with "; " as the separator. For single errors or
// unwrapped errors with one element, it returns the error's message directly.
// For non-joined errors, it returns the error's message as is.
func errorMessage(err error) string {
	if jes, ok := err.(joined); ok {
		es := jes.Unwrap()
		if len(es) == 1 {
			return es[0].Error()
		}
		b := []byte(es[0].Error())
		for _, err := range es[1:] {
			b = append(b, ';', ' ')
			b = append(b, err.Error()...)
		}
		// At this point, b has at least two bytes '\n' and ' '.
		return unsafe.String(&b[0], len(b)) // nolint: gosec
	}
	return err.Error()
}
