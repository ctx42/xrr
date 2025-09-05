package xrr

import (
	"time"
)

// meta represents metadata collection.
type meta struct {
	m map[string]any
}

// Meta returns a new metadata collection.
func Meta() meta { return meta{} }

// Str adds the key with string val to the metadata collection. Key will be
// overridden with the new value if it already exists.
func (m meta) Str(key, value string) meta {
	return m.set(key, value)
}

// Int adds the key with integer val to the metadata collection. Key will be
// overridden with the new value if it already exists.
func (m meta) Int(key string, value int) meta {
	return m.set(key, value)
}

// Int64 adds the key with int64 val to the metadata collection. Key will be
// overridden with a new value if it already exists.
func (m meta) Int64(key string, value int64) meta {
	return m.set(key, value)
}

// Float64 adds the key with float64 val to the metadata collection. Key will
// be overridden with the new value if it already exists.
func (m meta) Float64(key string, value float64) meta {
	return m.set(key, value)
}

// Bool adds the key with val as a boolean to the metadata collection. Key will be
// // overridden with the new value if it already exists.
func (m meta) Bool(key string, value bool) meta {
	return m.set(key, value)
}

// Time adds the key with val as a time to the metadata collection. Key will be
// // overridden with the new value if it already exists.
func (m meta) Time(key string, value time.Time) meta {
	return m.set(key, value)
}

// set sets instance metadata key/value if the metadata map is nil, it will
// allocate it.
func (m meta) set(key string, value any) meta {
	if m.m == nil {
		m.m = make(map[string]any)
	}
	m.m[key] = value
	return m
}
