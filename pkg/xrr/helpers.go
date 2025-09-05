package xrr

// DefaultCode returns the first non-empty code from the slice of codes.
func DefaultCode(otherwise string, codes ...string) string {
	for _, code := range codes {
		if code != "" {
			return code
		}
	}
	return otherwise
}
