// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac <rzajac@gmail.com>
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
		x, _ := assert.SameType(t, &GenericError[EDGeneric]{}, err)
		assert.Equal(t, "msg", x.Error())
		assert.Equal(t, "ECode", x.code)
		assert.Nil(t, x.meta)
	})

	t.Run("with options", func(t *testing.T) {
		// --- Given ---
		opt := Meta().Int("A", 1).Int("B", 2).Option()

		// --- When ---
		err := New("msg", "ECode", opt)

		// --- Then ---
		x, _ := assert.SameType(t, &GenericError[EDGeneric]{}, err)
		assert.Equal(t, "msg", x.Error())
		assert.Equal(t, "ECode", x.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, x.meta)
	})

	t.Run("WithCode overrides code argument", func(t *testing.T) {
		// --- Given ---
		opt := WithCode("MyCode")

		// --- When ---
		err := New("msg", "ECode", opt)

		// --- Then ---
		x, _ := assert.SameType(t, &GenericError[EDGeneric]{}, err)
		assert.Equal(t, "msg", x.Error())
		assert.Equal(t, "MyCode", x.code)
	})
}

func Test_Wrap(t *testing.T) {
	t.Run("wrapping nil returns nil", func(t *testing.T) {
		// --- When ---
		err := Wrap[string](nil)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("wrap error without options", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")

		// --- When ---
		err := Wrap[string](e)

		// --- Then ---
		assert.Same(t, e, errors.Unwrap(err))
	})

	t.Run("wrap std error and set error code", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		opt := WithCode("ECode")

		// --- When ---
		err := Wrap[string](e, opt)

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
		err := Wrap[string](e, opt)

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
		err := Wrap[string](e, opts...)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *GenericError[string]
		assert.Type(t, &x, err)
		assert.Same(t, e, x.Unwrap())
		assert.Equal(t, "ECode", x.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, x.meta)
	})
}
