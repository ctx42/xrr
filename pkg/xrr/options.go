// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

// Option represents an option for configuring [GenericError] instances.
type Option func(*Options)

// Options is a collection of options used in the xrr package.
type Options struct {
	code string         // Error code.
	meta map[string]any // Metadata associated with an error.
	err  error          // Wrapped error.
}

// Set applies the provided options to the [Options] instance and returns it.
func (ops Options) Set(opts ...Option) Options {
	for _, opt := range opts {
		opt(&ops)
	}
	return ops
}

// WithCode is an option for setting the error code.
func WithCode(code string) Option {
	return func(ops *Options) { ops.code = code }
}

// WithMeta is an option for setting the metadata. The provided map must not be
// modified or reused by the caller after passing it to this function. The
// value types that are not supported will be skipped.
//
// For supported metadata types see [MetaType] type constraint.
func WithMeta(meta map[string]any) Option {
	return func(ops *Options) {
		for key, value := range meta {
			if !isTypeSupported(value) {
				continue
			}
			if ops.meta == nil {
				ops.meta = make(map[string]any, len(meta))
			}
			ops.meta[key] = value
		}
	}
}

// WithMetaFrom is an option for [New] and [Wrap] setting the metadata from a
// [Metadater] instance. The types that are not supported by the [GenericError]
// metadata are not going to be added.
func WithMetaFrom(src Metadater) Option {
	return WithMeta(src.MetaAll())
}
