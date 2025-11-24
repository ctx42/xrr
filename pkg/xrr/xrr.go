// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package xrr provides errors supporting error codes and metadata.
package xrr

// ECGeneric represents generic error code used for non-nil errors which have
// no error code assigned.
const ECGeneric = "ECGeneric"

// Coder is the interface that wraps the ErrorCode method.
//
// ErrorCode returns error code.
//
// For nil errors it must return an empty string, but for non-nil errors
// without assigned code, it should return [ECGeneric] error code.
type Coder interface {
	ErrorCode() string
}

// Fielder is the interface that wraps the ErrorFields method.
//
// ErrorFields returns errors for field names. It is used for validation errors.
type Fielder interface {
	ErrorFields() map[string]error
}

// Metadater is an interface providing access to error metadata.
type Metadater interface {
	// MetaAll returns a copy of the metadata from the [Error] instance.
	//
	// It does not include metadata from wrapped errors. To retrieve metadata,
	// recursively, use the [GetMeta] function instead. Returns nil if no
	// metadata is present.
	MetaAll() map[string]any
}
