// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

// Sentinel error codes.
const (
	// ECInvJSON represents invalid JSON error code.
	ECInvJSON = "ECInvJSON"

	// ECInvJSONError represents the error code for a JSON string that has
	// invalid syntax or structure to serve as a [GenericError] representation.
	ECInvJSONError = "ECInvJSONError"

	// ECFields represents the [ErrFields] error code.
	ECFields = "ECFields"
)

// Sentinel errors.
var (
	// ErrInvJSON represents an error indicating JSON structure or format error.
	ErrInvJSON = New("invalid JSON", ECInvJSON)

	// ErrInvJSONError represents an error indicating a JSON string has invalid
	// syntax or structure to serve as a [GenericError] representation.
	ErrInvJSONError = New("invalid JSON error representation", ECInvJSONError)

	// ErrFields is the default lead error used by [Envelop] when the cause
	// implements [Fielder] and no explicit lead error is provided.
	ErrFields = New("fields error", ECFields)
)
