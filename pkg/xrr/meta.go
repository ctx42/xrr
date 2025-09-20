// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"time"
)

// metaType is a type of metadata value.
type metaType interface {
	bool | string | int | int64 | float64 | time.Time | time.Duration
}

// Metadata represents metadata collection.
type Metadata struct {
	m map[string]any
}

// Meta returns a new instance of [Metadata].
func Meta() Metadata { return Metadata{} }

// Bool adds the key with val as a boolean to the metadata collection. Key will be
// overridden with the new value if it already exists.
func (m Metadata) Bool(key string, value bool) Metadata {
	return m.set(key, value)
}

// Str adds the key with string val to the metadata collection. Key will be
// overridden with the new value if it already exists.
func (m Metadata) Str(key, value string) Metadata {
	return m.set(key, value)
}

// Int adds the key with integer val to the metadata collection. Key will be
// overridden with the new value if it already exists.
func (m Metadata) Int(key string, value int) Metadata {
	return m.set(key, value)
}

// Int64 adds the key with int64 val to the metadata collection. Key will be
// overridden with a new value if it already exists.
func (m Metadata) Int64(key string, value int64) Metadata {
	return m.set(key, value)
}

// Float64 adds the key with float64 val to the metadata collection. Key will
// be overridden with the new value if it already exists.
func (m Metadata) Float64(key string, value float64) Metadata {
	return m.set(key, value)
}

// Time adds the key with val as a time to the metadata collection. Key will be
// overridden with the new value if it already exists.
func (m Metadata) Time(key string, value time.Time) Metadata {
	return m.set(key, value)
}

// MetaSetAll copies all metadata from the given map. Only the supported types
// will be copied.
func (m Metadata) MetaSetAll(meta map[string]any) Metadata {
	for key, value := range meta {
		if m.m == nil {
			m.m = make(map[string]any)
		}
		if isTypeSupported(value) {
			m.m[key] = value
		}
	}
	return m
}

// MetaSetFrom copies all metadata from the given [Metadater] instance. Only
// the supported types will be copied.
func (m Metadata) MetaSetFrom(meta Metadater) Metadata {
	for key, value := range meta.MetaAll() {
		if m.m == nil {
			m.m = make(map[string]any)
		}
		if isTypeSupported(value) {
			m.m[key] = value
		}
	}
	return m
}

// Option returns a function that sets the metadata on the [Error] instance.
// TODO(rz): test not overriding existing metadata map.
func (m Metadata) Option() func(*Error) { return WithMeta(m.m) }

// set sets instance metadata key/value if the metadata map is nil, it will
// allocate it.
func (m Metadata) set(key string, value any) Metadata {
	if m.m == nil {
		m.m = make(map[string]any)
	}
	m.m[key] = value
	return m
}
