package xrr

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

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
