// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewErrors(t *testing.T) {
	// --- When ---
	ec := NewErrors()

	// --- Then ---
	assert.Nil(t, ec.First())
	assert.Len(t, 0, ec.Unwrap())
}

func Test_Errors_Add(t *testing.T) {
	// --- Given ---
	err0 := New("msg0", "ECode0")
	err1 := errors.New("msg1")
	ec := NewErrors()

	// --- When ---
	ec.Add(err0)
	ec.Add(err1)

	// --- Then ---
	es := ec.Unwrap()
	assert.Len(t, 2, es)
	assert.Same(t, err0, es[0])
	assert.Same(t, err1, es[1])
}

func Test_Errors_Add_First(t *testing.T) {
	// --- Given ---
	err0 := New("msg0", "ECode0")
	err1 := New("msg1", "ECode1")
	ec := NewErrors()

	// --- When ---
	ec.Add(err0)
	ec.Add(err1)

	// --- Then ---
	assert.Same(t, err0, ec.First())
	assert.Len(t, 2, ec.Unwrap())
}

func Test_Errors_Errors(t *testing.T) {
	// --- Given ---
	err0 := New("msg0", "ECode0")
	err1 := New("msg1", "ECode1")
	ec := NewErrors()

	// --- When ---
	ec.Add(err0)
	ec.Add(err1)

	// --- Then ---
	es := ec.Unwrap()
	assert.Len(t, 2, es)
	assert.Equal(t, err0, es[0])
	assert.Equal(t, err1, es[1])

	esi := ec.Unwrap()
	assert.Len(t, 2, esi)
	assert.Equal(t, err0, esi[0])
	assert.Equal(t, err1, esi[1])
}

func Test_Errors_Reset(t *testing.T) {
	// --- Given ---
	ec := NewErrors()
	ec.Add(New("msg0", "ECode0"))

	// --- When ---
	ec.Reset()

	// --- Then ---
	assert.Len(t, 0, ec.Unwrap())
	assert.Equal(t, 0, len(ec))
}
