// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
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
