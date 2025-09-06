package xrr

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithCode(t *testing.T) {
	// --- Given ---
	e := &Error{}

	// --- When ---
	WithCode("ECode")(e)

	// --- Then ---
	assert.Equal(t, "ECode", e.code)
}

func Test_New(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- When ---
		err := New("msg", "ECode")

		// --- Then ---
		assert.Equal(t, "msg", err.Error())
		assert.Equal(t, "ECode", err.code)
		assert.Nil(t, err.meta)
	})

	t.Run("with options", func(t *testing.T) {
		// --- Given ---
		opt := Meta().Int("A", 1).Int("B", 2).Option()

		// --- When ---
		err := New("msg", "ECode", opt)

		// --- Then ---
		assert.Equal(t, "msg", err.Error())
		assert.Equal(t, "ECode", err.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, err.meta)
	})

	t.Run("WithCode overrides code argument", func(t *testing.T) {
		// --- Given ---
		opt := WithCode("MyCode")

		// --- When ---
		err := New("msg", "ECode", opt)

		// --- Then ---
		assert.Equal(t, "msg", err.Error())
		assert.Equal(t, "MyCode", err.code)
	})
}

func Test_Error_ErrorCode(t *testing.T) {
	// --- Given ---
	err := &Error{code: "ECode"}

	// --- When ---
	have := err.ErrorCode()

	// --- Then ---
	assert.Equal(t, "ECode", have)
}

func Test_Error_MetaAll(t *testing.T) {
	t.Run("returns a clone", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{"A": 1, "B": 2}
		e := &Error{meta: m}

		// --- When ---
		have := e.MetaAll()

		// --- Then ---
		assert.NotSame(t, m, have)
		assert.Equal(t, m, have)
	})
}

func Test_Error_Unwrap(t *testing.T) {
	t.Run("returns wrapped error", func(t *testing.T) {
		// --- Given ---
		err := New("msg", "ECode")

		// --- When ---
		have := err.Unwrap()

		// --- Then ---
		assert.Same(t, err.error, have)
	})

	t.Run("returns nil for nil instance", func(t *testing.T) {
		// --- Given ---
		var err *Error

		// --- When ---
		have := err.Unwrap()

		// --- Then ---
		assert.Nil(t, have)
	})
}
