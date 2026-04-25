// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Options_Set(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- Given ---
		ops := Options{}

		// --- When ---
		have := ops.Set()

		// --- Then ---
		assert.Zero(t, have)
	})

	t.Run("with options", func(t *testing.T) {
		// --- Given ---
		ops := Options{}
		o0 := WithCode("ECode")
		o1 := WithMeta(map[string]any{"A": 1})

		// --- When ---
		have := ops.Set(o0, o1)

		// --- Then ---
		want := Options{
			code: "ECode",
			meta: map[string]any{"A": 1},
		}
		assert.Equal(t, want, have)
		assert.Zero(t, ops)
	})
}

func Test_WithCode(t *testing.T) {
	// --- Given ---
	ops := &Options{}

	// --- When ---
	WithCode("ECode")(ops)

	// --- Then ---
	assert.Equal(t, "ECode", ops.code)
}

func Test_WithMeta(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{
			"bool":          true,
			"string":        "abc",
			"int":           2,
			"int64":         int64(2),
			"float64":       4.2,
			"time":          time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration":      time.Second,
			"not-supported": func() {},
		}
		ops := &Options{}

		// --- When ---
		WithMeta(m)(ops)

		// --- Then ---
		want := map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		}
		assert.Equal(t, want, ops.meta)
	})

	t.Run("not supported types are removed", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{"A": 1, "B": struct{}{}}
		ops := &Options{}

		// --- When ---
		WithMeta(m)(ops)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, ops.meta)
	})

	t.Run("multiple calls work like merge", func(t *testing.T) {
		// --- Given ---
		m0 := map[string]any{"A": 1, "B": 2}
		m1 := map[string]any{"B": 3}
		ops := &Options{}

		// --- When ---
		WithMeta(m0)(ops)
		WithMeta(m1)(ops)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1, "B": 3}, ops.meta)
	})
}

func Test_WithCause(t *testing.T) {
	t.Run("sets err and inherits code when ops.code is empty", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECCause")
		ops := &Options{}

		// --- When ---
		WithCause(cause)(ops)

		// --- Then ---
		assert.Equal(t, cause, ops.err)
		assert.Equal(t, "ECCause", ops.code)
	})

	t.Run("it does not override existing code", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECCause")
		ops := &Options{code: "ECExisting"}

		// --- When ---
		WithCause(cause)(ops)

		// --- Then ---
		assert.Equal(t, cause, ops.err)
		assert.Equal(t, "ECExisting", ops.code)
	})
}

func Test_WithMetaFrom(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		m := TMetaAll(map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		})
		ops := &Options{}

		// --- When ---
		WithMetaFrom(m)(ops)

		// --- Then ---
		want := map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		}
		assert.Equal(t, want, ops.meta)
	})

	t.Run("not supported types are removed", func(t *testing.T) {
		// --- Given ---
		m := TMetaAll(map[string]any{"A": 1, "B": struct{}{}})
		ops := &Options{}

		// --- When ---
		WithMeta(m)(ops)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, ops.meta)
	})
}
