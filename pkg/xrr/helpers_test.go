package xrr

import (
	"errors"
	"testing"

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

func Test_prefix(t *testing.T) {
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
