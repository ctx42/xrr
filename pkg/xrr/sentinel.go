// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

const (
	// ECInvJSON represents invalid JSON error code.
	ECInvJSON = "ECInvJSON"

	// ECInvJSONError represents error code indicating a JSON string has
	// invalid syntax or structure to be the [Error] representation.
	ECInvJSONError = "ECInvJSONError"

	// ECFields represents generic [Fields] error code.
	ECFields = "ECFields"
)

var (
	// ErrInvJSON represents an error indicating JSON structure or format error.
	ErrInvJSON = New("invalid JSON", ECInvJSON)

	// ErrInvJSONError represents an error indicating a JSON string has invalid
	// syntax or structure to be the [Error] representation.
	ErrInvJSONError = New("invalid JSON error representation", ECInvJSONError)

	// ErrFields represents generic [Fields] error.
	ErrFields = New("fields error", ECFields)
)
