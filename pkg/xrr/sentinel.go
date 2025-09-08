package xrr

const (
	// ECInvJSON represents invalid JSON error code.
	ECInvJSON = "ECInvJSON"

	// ECInvJSONError represents error code indicating a JSON string has
	// invalid syntax or structure to be the [Error] representation.
	ECInvJSONError = "ECInvJSONError"
)

var (
	// ErrInvJSON represents an error indicating JSON structure or format error.
	ErrInvJSON = New("invalid JSON", ECInvJSON)

	// ErrInvJSONError represents an error indicating a JSON string has invalid
	// syntax or structure to be the [Error] representation.
	ErrInvJSONError = New("invalid JSON error representation", ECInvJSONError)
)
