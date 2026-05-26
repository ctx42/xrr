// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WrapUsing(t *testing.T) {
	t.Run("wrapping nil returns nil", func(t *testing.T) {
		// --- When ---
		err := WrapUsing[string](nil)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("wrapping typed nil returns nil", func(t *testing.T) {
		// --- Given ---
		var e *GenericError[EDXrr]

		// --- When ---
		err := WrapUsing[string](e)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("wrap error without options", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")

		// --- When ---
		err := WrapUsing[string](e)

		// --- Then ---
		assert.Same(t, e, errors.Unwrap(err))
	})

	t.Run("wrap std error and set error code", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		opt := WithCode("ECode")

		// --- When ---
		err := WrapUsing[string](e, opt)

		// --- Then ---
		assert.NotSame(t, e, err)
		assert.Same(t, e, errors.Unwrap(err))
		assert.Equal(t, "ECode", GetCode(err))
	})

	t.Run("wrap std error and add metadata", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		opt := Meta().Int("A", 1).Int("B", 2).Option()

		// --- When ---
		err := WrapUsing[string](e, opt)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Same(t, e, x.Unwrap())
		assert.Equal(t, ECGeneric, x.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, x.meta)
	})

	t.Run("wrap std error and add error code metadata", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		opts := []Option{
			Meta().Int("A", 1).Int("B", 2).Option(),
			WithCode("ECode"),
		}

		// --- When ---
		err := WrapUsing[string](e, opts...)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Same(t, e, x.Unwrap())
		assert.Equal(t, "ECode", x.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, x.meta)
	})

	t.Run("wrap std error and cause error", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		cause := errors.New("cause")
		opts := []Option{
			WithCause(cause),
		}

		// --- When ---
		err := WrapUsing[string](e, opts...)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Equal(t, "msg; cause", x.Error())
		assert.Equal(t, "ECGeneric", x.code)
		assert.ErrorIs(t, e, err)
		assert.ErrorIs(t, cause, err)
	})

	t.Run("wrap std error and cause errors chain", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		cause0 := errors.New("cause0")
		cause1 := fmt.Errorf("cause1: %w", cause0)
		cause2 := fmt.Errorf("cause2: %w", cause1)
		opts := []Option{
			WithCause(cause2),
		}

		// --- When ---
		err := WrapUsing[string](e, opts...)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Equal(t, "msg; cause2: cause1: cause0", x.Error())
		assert.Equal(t, "ECGeneric", x.code)
		assert.ErrorIs(t, e, err)
		assert.ErrorIs(t, cause0, err)
		assert.ErrorIs(t, cause1, err)
		assert.ErrorIs(t, cause2, err)
	})

	t.Run("inherits code from err not from cause", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECMsg")
		cause := New("cause", "ECCause")

		// --- When ---
		err := WrapUsing[string](e, WithCause(cause))

		// --- Then ---
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Equal(t, "ECMsg", x.code)
		assert.Equal(t, "msg; cause", x.Error())
		assert.ErrorIs(t, e, err)
		assert.ErrorIs(t, cause, err)
	})

	t.Run("cause code not inherited when err has no explicit code", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		cause := New("cause", "ECCause")

		// --- When ---
		err := WrapUsing[string](e, WithCause(cause))

		// --- Then ---
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Equal(t, "ECGeneric", x.code)
		assert.ErrorIs(t, cause, err)
	})

	t.Run("WithCode after WithCause overrides code", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECMsg")
		cause := New("cause", "ECCause")

		// --- When ---
		err := WrapUsing[string](e, WithCause(cause), WithCode("ECOverride"))

		// --- Then ---
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Equal(t, "ECOverride", x.code)
		assert.ErrorIs(t, e, err)
		assert.ErrorIs(t, cause, err)
	})

	t.Run("meta alongside cause", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		cause := errors.New("cause")
		meta := Meta().Int("A", 42).Option()

		// --- When ---
		err := WrapUsing[string](e, WithCause(cause), meta)

		// --- Then ---
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Equal(t, map[string]any{"A": 42}, x.meta)
		assert.Equal(t, "msg; cause", x.Error())
		assert.ErrorIs(t, e, err)
		assert.ErrorIs(t, cause, err)
	})

	t.Run("meta from cause chain", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		meta := Meta().Int("B", 10).Option()
		cause := New("cause", "ECCause", meta)

		// --- When ---
		err := WrapUsing[string](e, WithCause(cause))

		// --- Then ---
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Nil(t, x.meta)
		assert.Equal(t, map[string]any{"B": 10}, GetMeta(err))
	})

	t.Run("override cause meta on same key", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		meta0 := Meta().Int("A", 1).Int("B", 10).Option()
		cause := New("cause", "ECCause", meta0)

		// --- When ---
		meta1 := Meta().Int("A", 99).Option()
		err := WrapUsing[string](e, WithCause(cause), meta1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 99, "B": 10}, GetMeta(err))
	})
}

func Test_SetCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		e := errors.New("error")

		// --- When ---
		err := SetCode[EDXrr](e, "ECode")

		// --- Then ---
		var xe *GenericError[EDXrr]
		assert.Type(t, &xe, err)
		assert.Same(t, e, xe.Unwrap())
		assert.Equal(t, "ECode", xe.ErrorCode())
	})

	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		err := SetCode[EDXrr](nil, "ECode")

		// --- Then ---
		assert.NoError(t, err)
	})

	t.Run("it does not wrap when the code is the same", func(t *testing.T) {
		// --- Given ---
		e := New("error", "ECode")

		// --- When ---
		err := SetCode[EDXrr](e, "ECode")

		// --- Then ---
		assert.Same(t, e, err)
	})

	t.Run("returns the same instance when code is empty", func(t *testing.T) {
		// --- Given ---
		e := errors.New("error")

		// --- When ---
		err := SetCode[EDXrr](e, "")

		// --- Then ---
		assert.Same(t, e, err)
	})
}
