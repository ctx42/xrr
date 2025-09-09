package xrr

import (
	"encoding/json"
	"fmt"
	"maps"
)

// WithCode is an option for [New] and [Wrap] setting the error code.
func WithCode(code string) func(*Error) {
	return func(e *Error) { e.code = code }
}

// Error represents an error with an error code and structured metadata.
type Error struct {
	msg  string         // Error message.
	code string         // Error code.
	meta map[string]any // Structured metadata.
	err  error          // Wrapped error.
}

// New creates a new [Error] instance with the specified message and error code.
// If the [WithCode] option is on the list of options, it will override the
// code argument.
func New(msg, code string, opts ...func(*Error)) error {
	err := &Error{
		msg:  msg,
		code: code,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

// Wrap wraps an error in an [Error] instance, applying the given options.
//
// It returns nil if the input error is nil or no options were provided. The
// returned error retains the same error code as the input error, obtained via
// [GetCode] function. To override the error code, use the [WithCode] option.
func Wrap(err error, opts ...func(*Error)) error {
	if err == nil {
		return nil
	}
	if len(opts) == 0 {
		return err
	}
	e := &Error{err: err, code: GetCode(err)}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Error) Error() string {
	if e.msg == "" && e.err != nil {
		return e.err.Error()
	}
	return e.msg
}

// ErrorCode returns error code.
func (e *Error) ErrorCode() string { return e.code }

// MetaAll returns a clone of the error's metadata.
func (e *Error) MetaAll() map[string]any { return maps.Clone(e.meta) }

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *Error) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"error": e.Error(),
		"code":  e.code,
	}
	if len(e.meta) > 0 {
		m["meta"] = e.meta
	}
	return json.Marshal(m)
}

// UnmarshalJSON unmarshal JSON representation of the [Error].
//
// The minimal valid JSON representation for an [Error] is
//
//	{"error": "message"}
//
// and in this case, the error code is set to [ECGeneric].
//
// Notes:
//   - all metadata numeric values will be unmarshalled as float64
func (e *Error) UnmarshalJSON(data []byte) error {
	m := make(map[string]any, 3)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	msgI, _ := m["error"]
	msg, _ := msgI.(string)
	if msg == "" {
		return ErrInvJSONError
	}

	codeI, _ := m["code"]
	code, _ := codeI.(string)
	if code == "" {
		code = ECGeneric
	}

	metaI, _ := m["meta"]
	var meta map[string]any
	if metaI != nil {
		meta, _ = metaI.(map[string]any)
	}

	e.msg = msg
	e.code = code
	e.meta = meta
	return nil
}

func (e *Error) Format(state fmt.State, verb rune) {
	Format(e.Error(), e.code, state, verb)
}

// Format is a custom formatter for Immutable and Error instances.
func Format(msg, code string, state fmt.State, verb rune) {
	switch verb {
	case 's', 'q':
		if verb == 'q' {
			msg = fmt.Sprintf("%q", msg)
		}
		_, _ = fmt.Fprint(state, msg)

	case 'v':
		_, _ = fmt.Fprint(state, msg)
		if state.Flag('+') {
			_, _ = fmt.Fprintf(state, " (%s)", code)
		}
	}
}
