package xrr

import (
	"errors"
)

// GetCode returns error code. If the error implements [Coder] interface, it
// will be used. Otherwise, it goes through error chains (and joins) to find
// the first non-empty error code.
//
// For nil error it will return an empty string. If the error code is not found,
// returns [ECGeneric] error code.
func GetCode(err error) string {
	if err == nil || isNil(err) {
		return ""
	}

	switch e := err.(type) { // nolint: errorlint
	case Coder:
		return e.ErrorCode()

	case interface{ Unwrap() []error }:
		// Get the first code from the joined errors.
		for _, je := range e.Unwrap() {
			return GetCode(je)
		}
	}

	// We really go out of our way to get an error code.
	if c := GetCode(errors.Unwrap(err)); c != "" {
		return c
	}
	return ECGeneric
}

// GetMeta recursively retrieves metadata from an error and its wrapped errors.
//
// It returns a map containing metadata from the error, merging metadata in a
// top-down order (metadata from the outermost error takes precedence). If any
// of the errors implements the MetaAll method, it's called to get its metadata.
// For errors implementing Unwrap() []error, metadata is collected from all
// wrapped errors, prioritizing earlier errors in the sequence over later ones.
// Returns nil if no metadata is found or the error is nil.
func GetMeta(err error) map[string]any {
	if err == nil || isNil(err) {
		return nil
	}

	switch e := err.(type) { // nolint: errorlint
	case Metadater:
		hi := e.MetaAll()
		lo := GetMeta(errors.Unwrap(err))
		if len(lo) == 0 {
			return hi
		}
		for k, v := range hi {
			lo[k] = v
		}
		return lo

	case interface{ Unwrap() []error }:
		var meta map[string]any
		ers := e.Unwrap()
		for i := len(ers) - 1; i >= 0; i-- {
			if m := GetMeta(ers[i]); m != nil {
				if meta == nil {
					meta = make(map[string]any, len(m))
				}
				for k, v := range m {
					meta[k] = v
				}
			}
		}
		return meta
	}
	return GetMeta(errors.Unwrap(err))
}
