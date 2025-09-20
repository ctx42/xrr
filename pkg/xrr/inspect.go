// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"time"
)

// IsCode walks the error chain (tree) and returns true if any of the errors
// has a given error code.
func IsCode(err error, code string) bool {
	var is bool
	cb := func(err error) bool {
		if GetCode(err) == code {
			is = true
			return false
		}
		return true
	}
	walk(err, cb)
	return is
}

// GetCode returns error code associated with the provided error. If an error
// does not implement [Coder] interface, the [ECGeneric] error code is returned.
// For nil error it will return an empty string.
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
// and its wrapped errors, ignoring empty error codes.
func GetCodes(err error) []string {
	set := make(map[string]struct{}, 10)
	var ret []string
	cb := func(err error) bool {
		code := GetCode(err)
		if _, ok := set[code]; !ok {
			set[code] = struct{}{}
			ret = append(ret, code)
		}
		return true
	}
	walk(err, cb)
	return ret
}

// GetMeta recursively retrieves metadata from an error and its wrapped errors.
//
// The error chain (tree) is traversed using the breath-first search approach
// with errors closer to the top and more on the left override metadata from
// the lover and more to the right parts of the tree.
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
	walkReverse(err, cb)
	return m
}

// GetBool recursively walks the error chain (tree) and returns the first bool
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns a false and false.
func GetBool(err error, key string) (bool, bool) {
	return getKey[bool](err, key)
}

// GetStr recursively walks the error chain (tree) and returns the first string
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns an empty string and false.
func GetStr(err error, key string) (string, bool) {
	return getKey[string](err, key)
}

// GetInt recursively walks the error chain (tree) and returns the first int
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns a zero value and false.
func GetInt(err error, key string) (int, bool) {
	return getKey[int](err, key)
}

// GetInt64 recursively walks the error chain (tree) and returns the first
// int64 value associated with the provided key. Returns the key value and true
// if the key was found. Otherwise, returns a zero value and false.
func GetInt64(err error, key string) (int64, bool) {
	return getKey[int64](err, key)
}

// GetFloat64 recursively walks the error chain (tree) and returns the first
// float64 value associated with the provided key. Returns the key value and
// true if the key was found. Otherwise, returns a zero value and false.
func GetFloat64(err error, key string) (float64, bool) {
	return getKey[float64](err, key)
}

// GetTime recursively walks the error chain (tree) and returns the first
// [time.Time] value associated with the provided key. Returns the key value
// and true if the key was found. Otherwise, returns a zero value and false.
func GetTime(err error, key string) (time.Time, bool) {
	return getKey[time.Time](err, key)
}

// GetDuration recursively walks the error chain (tree) and returns the first
// [time.Duration] value associated with the provided key. Returns the key value
// and true if the key was found. Otherwise, returns a zero value and false.
func GetDuration(err error, key string) (time.Duration, bool) {
	return getKey[time.Duration](err, key)
}

// getKey recursively walks the error chain (tree) and returns the first string
// value associated with the provided key. Returns the key value and true if
// the key was found. Otherwise, returns an empty string and false.
func getKey[T metaType](err error, key string) (T, bool) {
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
	walk(err, cb)
	return value, found
}

// walk walks the error chain (tree) using breadth-first search (BFS) and calls
// the callback for each error. Return true from the callback if you want to
// continue walking the tree or false to stop.
func walk(err error, cb func(err error) bool) bool {
	if err == nil || isNil(err) {
		return true
	}
	switch x := err.(type) { // nolint: errorlint
	case interface{ Unwrap() error }:
		if !cb(err) {
			return false
		}
		if e := x.Unwrap(); e != nil {
			return walk(e, cb)
		}
		return true

	case Fielder:
		_, ers := sortFields(x.ErrorFields())
		for _, fe := range ers {
			if !walk(fe, cb) {
				return false
			}
		}
		return true

	case joined:
		for _, je := range x.Unwrap() {
			if !walk(je, cb) {
				return false
			}
		}
		return true
	}
	return cb(err)
}

// walkReverse works like [walk] but in the reverse order.
func walkReverse(err error, cb func(err error) bool) bool {
	if err == nil || isNil(err) {
		return true
	}
	switch x := err.(type) { // nolint: errorlint
	case interface{ Unwrap() error }:
		if e := x.Unwrap(); e != nil {
			if !walkReverse(e, cb) {
				return false
			}
		}

	case Fielder:
		_, ers := sortFields(x.ErrorFields())
		for i := len(ers) - 1; i >= 0; i-- {
			if !walkReverse(ers[i], cb) {
				return false
			}
		}
		return true

	case joined:
		ers := x.Unwrap()
		for i := len(ers) - 1; i >= 0; i-- {
			if !walkReverse(ers[i], cb) {
				return false
			}
		}
		return true
	}
	return cb(err)
}
