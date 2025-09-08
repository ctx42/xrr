package xrr

import (
	"encoding/json"
	"errors"
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
		var x *Error
		assert.Type(t, &x, err)
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
		var x *Error
		assert.Type(t, &x, err)
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
		var x *Error
		assert.Type(t, &x, err)
		assert.Equal(t, "msg", x.Error())
		assert.Equal(t, "MyCode", x.code)
	})
}

func Test_Wrap(t *testing.T) {
	t.Run("wrapping nil returns nil", func(t *testing.T) {
		// --- When ---
		err := Wrap(nil)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("wrap error without options", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")

		// --- When ---
		err := Wrap(e)

		// --- Then ---
		assert.Same(t, e, err)
	})

	t.Run("wrap std error and set error code", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		opt := WithCode("ECode")

		// --- When ---
		err := Wrap(e, opt)

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
		err := Wrap(e, opt)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *Error
		assert.Type(t, &x, err)
		assert.Same(t, e, x.Unwrap())
		assert.Equal(t, ECGeneric, x.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, x.meta)
	})

	t.Run("wrap std error and add error code metadata", func(t *testing.T) {
		// --- Given ---
		e := errors.New("msg")
		opts := []func(*Error){
			Meta().Int("A", 1).Int("B", 2).Option(),
			WithCode("ECode"),
		}

		// --- When ---
		err := Wrap(e, opts...)

		// --- Then ---
		assert.NotSame(t, e, err)
		var x *Error
		assert.Type(t, &x, err)
		assert.Same(t, e, x.Unwrap())
		assert.Equal(t, "ECode", x.code)
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, x.meta)
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
		have := errors.Unwrap(err)

		// --- Then ---
		var x *Error
		assert.Type(t, &x, err)
		assert.Same(t, x.error, have)
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

func Test_Error_MarshalJSON(t *testing.T) {
	t.Run("without metadata", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECode")

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		exp := `{"error":"msg", "code":"ECode"}`
		assert.JSON(t, exp, string(data))
	})

	t.Run("with metadata", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECode", Meta().Str("key", "val").Option())

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		exp := `{"error":"msg", "code":"ECode", "meta": {"key": "val"}}`
		assert.JSON(t, exp, string(data))
	})
}
