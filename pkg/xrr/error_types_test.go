// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_New(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- When ---
		err := New("msg", "ECode")

		// --- Then ---
		e, _ := assert.SameType(t, &GenericError[EDXrr]{}, err)
		assert.Equal(t, "msg", e.Error())
		assert.Equal(t, "ECode", e.code)
		assert.Nil(t, e.meta)
	})

	t.Run("with options", func(t *testing.T) {
		// --- Given ---
		opt := Meta().Int("A", 1).Int("B", 2).Option()

		// --- When ---
		err := New("msg", "ECode", opt)

		// --- Then ---
		e, _ := assert.SameType(t, &GenericError[EDXrr]{}, err)
		assert.Equal(t, "msg", e.Error())
		assert.Equal(t, "ECode", e.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, e.meta)
	})

	t.Run("WithCode overrides code argument", func(t *testing.T) {
		// --- Given ---
		opt := WithCode("MyCode")

		// --- When ---
		err := New("msg", "ECode", opt)

		// --- Then ---
		e, _ := assert.SameType(t, &GenericError[EDXrr]{}, err)
		assert.Equal(t, "msg", e.Error())
		assert.Equal(t, "MyCode", e.code)
	})
}

func Test_NewFieldError(t *testing.T) {
	t.Run("not nil error", func(t *testing.T) {
		// --- Given ---
		err := errors.New("msg")

		// --- When ---
		have := NewFieldError("name", err)

		// --- Then ---
		e, _ := assert.SameType(t, &GenericFields[EDXrr]{}, have)
		assert.Equal(t, 1, e.Len())
		assert.ErrorEqual(t, "name: msg", have)
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		have := NewFieldError("name", nil)

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_NewFieldErrors(t *testing.T) {
	t.Run("creates EDGeneric fields from map", func(t *testing.T) {
		// --- Given ---
		m := map[string]error{"f0": errors.New("em0")}

		// --- When ---
		have := NewFieldErrors(m)

		// --- Then ---
		fs, _ := assert.SameType(t, &GenericFields[EDXrr]{}, have)
		assert.Equal(t, 1, fs.Len())
		assert.ErrorEqual(t, "em0", fs.fields["f0"])
	})

	t.Run("map is stored directly without copying", func(t *testing.T) {
		// --- Given ---
		m := map[string]error{"f0": errors.New("em0")}

		// --- When ---
		have := NewFieldErrors(m)

		// --- Then ---
		fs, _ := assert.SameType(t, &GenericFields[EDXrr]{}, have)
		m["f1"] = errors.New("em1")
		assert.Equal(t, 2, fs.Len())
	})
}

func Test_Wrap(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		// --- When ---
		have := Wrap(nil)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("wraps error preserving message", func(t *testing.T) {
		// --- Given ---
		cause := errors.New("cause msg")

		// --- When ---
		have := Wrap(cause)

		// --- Then ---
		e, _ := assert.SameType(t, &Error{}, have)
		assert.Equal(t, "cause msg", e.Error())
		assert.ErrorIs(t, cause, have)
	})

	t.Run("inherits code from cause", func(t *testing.T) {
		// --- Given ---
		cause := New("msg", "CauseCode")

		// --- When ---
		have := Wrap(cause)

		// --- Then ---
		e, _ := assert.SameType(t, &Error{}, have)
		assert.Equal(t, "CauseCode", e.code)
	})

	t.Run("WithCode overrides inherited code", func(t *testing.T) {
		// --- Given ---
		cause := New("msg", "CauseCode")
		opt := WithCode("NewCode")

		// --- When ---
		have := Wrap(cause, opt)

		// --- Then ---
		x, _ := assert.SameType(t, &Error{}, have)
		assert.Equal(t, "NewCode", x.code)
	})

	t.Run("with metadata", func(t *testing.T) {
		// --- Given ---
		cause := errors.New("cause msg")
		opt := Meta().Str("key", "val").Option()

		// --- When ---
		have := Wrap(cause, opt)

		// --- Then ---
		e, _ := assert.SameType(t, &Error{}, have)
		assert.Equal(t, map[string]any{"key": "val"}, e.meta)
	})
}
