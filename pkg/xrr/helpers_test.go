// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Split(t *testing.T) {
	t.Run("joined error", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("msg0")
		e1 := New("msg1", "ECode")
		ers := errors.Join(e0, e1)

		// --- When ---
		have := Split(ers)

		// --- Then ---
		assert.Len(t, 2, have)
		assert.Same(t, e0, have[0])
		assert.Same(t, e1, have[1])
	})

	t.Run("single error", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("msg0")
		ers := errors.Join(e0)

		// --- When ---
		have := Split(ers)

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Same(t, e0, have[0])
	})

	t.Run("not joined error", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("msg0")

		// --- When ---
		have := Split(e0)

		// --- Then ---
		assert.Len(t, 1, have)
		assert.Same(t, e0, have[0])
	})

	t.Run("nil error", func(t *testing.T) {
		// --- Given ---
		var err error

		// --- When ---
		have := Split(err)

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_Join(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		// --- Given ---
		ers := []error{nil, nil, nil}

		// --- When ---
		have := Join(ers...)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("one error", func(t *testing.T) {
		// --- Given ---
		e1 := errors.New("m1")

		// --- When ---
		have := Join(e1)

		// --- Then ---
		assert.Same(t, e1, have)
	})

	t.Run("no gaps", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("m0")
		e1 := errors.New("m1")
		e2 := errors.New("m2")
		ers := []error{e0, e1, e2}

		// --- When ---
		have := Join(ers...)

		// --- Then ---
		assert.ErrorEqual(t, "m0\nm1\nm2", have)
	})
}

func Test_join(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		// --- Given ---
		var ers []error

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("all nil", func(t *testing.T) {
		// --- Given ---
		ers := []error{nil, nil, nil}

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("no gaps", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("m0")
		e1 := errors.New("m1")
		e2 := errors.New("m2")
		ers := []error{e0, e1, e2}

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Len(t, 3, have)
		assert.Same(t, e0, have[0])
		assert.Same(t, e1, have[1])
		assert.Same(t, e2, have[2])
	})

	t.Run("gap at the start", func(t *testing.T) {
		// --- Given ---
		e1 := errors.New("m0")
		e2 := errors.New("m2")
		ers := []error{nil, e1, e2}

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Len(t, 2, have)
		assert.Same(t, e1, have[0])
		assert.Same(t, e2, have[1])
	})

	t.Run("gap in the middle", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("m0")
		e2 := errors.New("m2")
		ers := []error{e0, nil, e2}

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Len(t, 2, have)
		assert.Same(t, e0, have[0])
		assert.Same(t, e2, have[1])
	})

	t.Run("gap at the end", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("m0")
		e1 := errors.New("m2")
		ers := []error{e0, e1, nil}

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Len(t, 2, have)
		assert.Same(t, e0, have[0])
		assert.Same(t, e1, have[1])
	})

	t.Run("gaps", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("m0")
		e1 := errors.New("m1")
		e2 := errors.New("m2")
		e3 := errors.New("m3")
		ers := []error{nil, e0, nil, nil, e1, nil, e2, nil, nil, e3}

		// --- When ---
		have := join(ers...)

		// --- Then ---
		assert.Len(t, 4, have)
		assert.Same(t, e0, have[0])
		assert.Same(t, e1, have[1])
		assert.Same(t, e2, have[2])
		assert.Same(t, e3, have[3])
	})
}

func Test_IsJoined(t *testing.T) {
	t.Run("joined error", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("msg0")
		e1 := New("msg1", "ECode")
		ers := errors.Join(e0, e1)

		// --- When ---
		have := IsJoined(ers)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not joined error", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("msg0")

		// --- When ---
		have := IsJoined(e0)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("nil error", func(t *testing.T) {
		// --- Given ---
		var err error

		// --- When ---
		have := IsJoined(err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_DefaultCode(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		// --- When ---
		have := DefaultCode("ECode")

		// --- Then ---
		assert.Equal(t, "ECode", have)
	})

	t.Run("empty slice", func(t *testing.T) {
		// --- When ---
		have := DefaultCode("ECode", []string{}...)

		// --- Then ---
		assert.Equal(t, "ECode", have)
	})

	t.Run("the first", func(t *testing.T) {
		// --- When ---
		have := DefaultCode("ECode", "First", "Second", "Third")

		// --- Then ---
		assert.Equal(t, "First", have)
	})

	t.Run("the first non-empty", func(t *testing.T) {
		// --- When ---
		have := DefaultCode("ECode", "", "First", "Second", "Third")

		// --- Then ---
		assert.Equal(t, "First", have)
	})
}

func Test_isNil_tabular(t *testing.T) {
	var err error

	tt := []struct {
		testN string

		value any
		want  bool
	}{
		{"nil", nil, true},
		{"typed nil", err, true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := isNil(tc.value)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_prefix_tabular(t *testing.T) {
	tt := []struct {
		testN string

		prefix string
		key    string
		exp    string
	}{
		{"1", "", "key", "key"},
		{"2", "pref", "key", "pref.key"},
		{"3", "pref", "", "pref"},
		{"4", "", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := prefix(tc.prefix, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exp, have)
		})
	}
}

func Test_isTypeSupported_tabular(t *testing.T) {
	tt := []struct {
		testN string

		typ  any
		want bool
	}{
		{"bool", true, true},
		{"string", "abc", true},
		{"int", 42, true},
		{"int64", int64(42), true},
		{"float64", 4.2, true},
		{"time", time.Now(), true},
		{"duration", time.Second, true},
		{"not supported", struct{}{}, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := isTypeSupported(tc.typ)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_sortFields(t *testing.T) {
	// --- Given ---
	fs := Fields{
		"f0": errors.New("em0"),
		"f1": nil,
		"f2": errors.New("em2"),
		"f3": nil,
		"f4": errors.New("em4"),
		"f5": nil,
	}

	// --- When ---
	hFields, hErs := sortFields(fs)

	// --- Then ---
	assert.Len(t, 6, hFields)
	assert.Equal(t, []string{"f0", "f1", "f2", "f3", "f4", "f5"}, hFields)

	assert.ErrorEqual(t, "em0", hErs[0])
	assert.Nil(t, hErs[1])
	assert.ErrorEqual(t, "em2", hErs[2])
	assert.Nil(t, hErs[3])
	assert.ErrorEqual(t, "em4", hErs[4])
	assert.Nil(t, hErs[5])
}
