package xrr

// Errors represent a collection of errors.
type Errors []error

// NewErrors returns a new instance of [Errors].
func NewErrors() Errors { return Errors{} }

// Add adds error to the collection.
func (ec *Errors) Add(e error) { *ec = append(*(ec), e) }

// Unwrap returns collected errors (MUST be treated as read-only).
func (ec *Errors) Unwrap() []error { return *ec }

// Reset resets the error collection.
func (ec *Errors) Reset() { *ec = (*ec)[:0] }

// First returns the first error in the collection or nil.
func (ec *Errors) First() error {
	if len(*ec) > 0 {
		return (*ec)[0]
	}
	return nil
}
