// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"time"
)

// IsCode walks the error chain (tree) using [walkErrors] and returns true if
// any error in the chain has the given code.
func IsCode(err error, code string) bool {
	var is bool
	cb := func(err error) bool {
		if GetCode(err) == code {
			is = true
			return false
		}
		return true
	}
	walkErrors(err, cb)
	return is
}

// GetCode returns the error code associated with the provided error. If an
// error does not implement the [Coder] interface, the [ECGeneric] error code
// is returned. For a nil error it returns an empty string.
func GetCode(err error) string {
	if err == nil || isNil(err) {
		return ""
	}
	if e, ok := err.(Coder); ok {
		return e.ErrorCode()
	}
	return ECGeneric
}

// GetCodes recursively retrieves a unique list of error codes from an error
// and its wrapped errors (via [walkErrors]), ignoring empty codes.
func GetCodes(err error) []string {
	set := make(map[string]struct{})
	var ret []string
	cb := func(err error) bool {
		code := GetCode(err)
		if code == "" {
			return true
		}
		if _, ok := set[code]; !ok {
			set[code] = struct{}{}
			ret = append(ret, code)
		}
		return true
	}
	walkErrors(err, cb)
	return ret
}

// GetMeta recursively retrieves metadata from an error and its wrapped errors
// (using [walkErrorsReverse] so that metadata closer to the root takes
// precedence).
func GetMeta(err error) map[string]any {
	var m map[string]any
	cb := func(err error) bool {
		if e, ok := err.(Metadater); ok {
			if meta := e.MetaAll(); len(meta) > 0 {
				if m == nil {
					m = make(map[string]any, len(meta))
				}
				for k, v := range meta {
					m[k] = v
				}
			}
		}
		return true
	}
	walkErrorsReverse(err, cb)
	return m
}

// GetBool recursively walks the error chain (tree) and returns the first bool
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns false and false.
func GetBool(err error, key string) (bool, bool) {
	return getMeta[bool](err, key)
}

// GetStr recursively walks the error chain (tree) and returns the first string
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns an empty string and false.
func GetStr(err error, key string) (string, bool) {
	return getMeta[string](err, key)
}

// GetInt recursively walks the error chain (tree) and returns the first int
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns a zero value and false.
func GetInt(err error, key string) (int, bool) {
	return getMeta[int](err, key)
}

// GetInt64 recursively walks the error chain (tree) and returns the first
// int64 value associated with the provided key. Returns the key value and true
// if the key was found. Otherwise, returns a zero value and false.
func GetInt64(err error, key string) (int64, bool) {
	return getMeta[int64](err, key)
}

// GetFloat64 recursively walks the error chain (tree) and returns the first
// float64 value associated with the provided key. Returns the key value and
// true if the key was found. Otherwise, returns a zero value and false.
func GetFloat64(err error, key string) (float64, bool) {
	return getMeta[float64](err, key)
}

// GetTime recursively walks the error chain (tree) and returns the first
// [time.Time] value associated with the provided key. Returns the key value
// and true if the key was found. Otherwise, returns a zero value and false.
func GetTime(err error, key string) (time.Time, bool) {
	return getMeta[time.Time](err, key)
}

// GetDuration recursively walks the error chain (tree) and returns the first
// [time.Duration] value associated with the provided key. Returns the key value
// and true if the key was found. Otherwise, returns a zero value and false.
func GetDuration(err error, key string) (time.Duration, bool) {
	return getMeta[time.Duration](err, key)
}

// getMeta walks the error chain using [walkErrors] and returns the
// first metadata value of type T associated with the given key.
//
// Traversal is breadth-first, so the first match encountered wins.
// Returns the zero value + false if the key is not found or the stored
// value has a different type.
//
// This is the shared implementation behind the typed Get* functions
// (GetStr, GetInt, GetBool, GetInt64, GetFloat64, GetTime, GetDuration).
func getMeta[T MetaType](err error, key string) (T, bool) {
	var value T
	var found bool
	cb := func(err error) bool {
		if e, ok := err.(Metadater); ok {
			if meta := e.MetaAll(); len(meta) > 0 {
				if v, exist := meta[key]; exist {
					if vv, success := v.(T); success {
						value = vv
						found = true
						return false
					}
				}
			}
		}
		return true
	}
	walkErrors(err, cb)
	return value, found
}

// walkErrors traverses an error chain/tree in breadth-first order,
// invoking the provided callback for each visited error.
//
// It understands three kinds of error structures:
//   - Plain wrapped errors (via Unwrap() error)
//   - Field errors (via [Fielder]), but only those that also implement [Coder]
//     are visited as nodes (plain field maps are treated as transparent)
//   - Joined errors (via Unwrap() []error)
//
// The callback should return true to continue traversal or false to stop early.
// This function is used internally to implement the various Get* functions.
//
// nolint: cyclop
func walkErrors(err error, cb func(err error) bool) bool {
	if err == nil || isNil(err) {
		return true
	}

	// We handle three shapes of error "trees":
	//   1. Single-wrapped errors (Unwrap() error)
	//   2. Field errors (Fielder) — only the ones that also
	//      implement Coder are visited
	//   3. Joined errors (Unwrap() []error)
	switch x := err.(type) { // nolint: errorlint
	case interface{ Unwrap() error }:
		if !cb(err) {
			return false
		}
		if e := x.Unwrap(); e != nil {
			return walkErrors(e, cb)
		}
		return true

	case Fielder:
		// Only visit [Fielder] nodes that also implement [Coder]; plain
		// field-map types (e.g. [GenericFields]) are transparent containers.
		if _, ok := err.(Coder); ok {
			if !cb(err) {
				return false
			}
		}
		_, ers := sortFieldErrors(x.ErrorFields())
		for _, fe := range ers {
			if !walkErrors(fe, cb) {
				return false
			}
		}
		return true

	case joined:
		for _, je := range x.Unwrap() {
			if !walkErrors(je, cb) {
				return false
			}
		}
		return true
	}

	// Plain error that doesn't implement any of the special interfaces above.
	return cb(err)
}

// walkErrorsReverse traverses an error chain/tree in reverse depth-first
// order (children before parents).
//
// It is primarily used by [GetMeta] so that metadata attached to errors
// closer to the root of the chain takes precedence over metadata from deeper
// errors.
//
// The three supported error structures and callback semantics are the same
// as [walkErrors].
func walkErrorsReverse(err error, cb func(err error) bool) bool {
	if err == nil || isNil(err) {
		return true
	}

	// Same three error shapes as walkErrors(), but we always visit children before
	// the parent (depth-first, reverse order). This ensures metadata from
	// outer errors wins when GetMeta walks in reverse.
	switch x := err.(type) { // nolint: errorlint
	case interface{ Unwrap() error }:
		// Recurse first (children before parent) so that deeper errors are
		// visited before this one.
		if e := x.Unwrap(); e != nil {
			if !walkErrorsReverse(e, cb) {
				return false
			}
		}

	case Fielder:
		_, ers := sortFieldErrors(x.ErrorFields())
		for i := len(ers) - 1; i >= 0; i-- {
			if !walkErrorsReverse(ers[i], cb) {
				return false
			}
		}
		// Only visit [Fielder] nodes that also implement [Coder]; plain
		// field-map types (e.g. [GenericFields]) are transparent containers.
		if _, ok := err.(Coder); ok {
			return cb(err)
		}
		return true

	case joined:
		ers := x.Unwrap()
		for i := len(ers) - 1; i >= 0; i-- {
			if !walkErrorsReverse(ers[i], cb) {
				return false
			}
		}
		return true
	}

	// Plain error (visited after its children, if any).
	return cb(err)
}
