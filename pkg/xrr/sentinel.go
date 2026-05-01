// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

// Sentinel error codes.
const (
	// ECInvJSON represents invalid JSON error code.
	ECInvJSON = "ECInvJSON"

	// ECInvJSONError represents error code indicating a JSON string has
	// invalid syntax or structure to be the [GenericError] representation.
	ECInvJSONError = "ECInvJSONError"

	// ECFields represents the [ErrFields] error code.
	ECFields = "ECFields"
)

// Sentinel errors.
var (
	// ErrInvJSON represents an error indicating JSON structure or format error.
	ErrInvJSON = New("invalid JSON", ECInvJSON)

	// ErrInvJSONError represents an error indicating a JSON string has invalid
	// syntax or structure to be the [GenericError] representation.
	ErrInvJSONError = New("invalid JSON error representation", ECInvJSONError)

	// ErrFields is the default lead error used by [Enclose] when the cause
	// implements [Fielder] and no explicit lead error is provided.
	ErrFields = New("fields error", ECFields)
)
