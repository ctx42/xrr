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

// Fielder is the interface that wraps the ErrFields method.
//
// ErrFields returns a list of errors for field names. It is mostly useful when
// returning validation errors.
type Fielder interface {
	ErrFields() map[string]error
}

// CodeFielder combines [Coder] and [Fielder] interfaces.
type CodeFielder interface {
	Coder
	Fielder
}
