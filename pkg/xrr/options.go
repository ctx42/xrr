// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
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

// WithMetaFrom is an option for [New] and [WrapUsing] setting the metadata
// from a [Metadater] instance. The types that are not supported by the
// [GenericError] metadata are not going to be added.
func WithMetaFrom(src Metadater) Option {
	return WithMeta(src.MetaAll())
}

// WithCause is an option for setting the wrapped cause error. The cause is
// accessible via [errors.Unwrap] and participates in [errors.Is] /
// [errors.As] chain traversal.
//
// If the code field is empty at the time this option is applied, it is
// inherited from cause via [GetCode]. Because options are applied in order,
// placing [WithCode] before [WithCause] causes [WithCause] to overwrite the
// code with the inherited value. Place [WithCode] after [WithCause] (or pass
// the code as the positional argument to [New]) to ensure the explicit code
// wins.
func WithCause(cause error) Option {
	return func(ops *Options) {
		if ops.code == "" {
			ops.code = GetCode(cause)
		}
		ops.err = cause
	}
}
