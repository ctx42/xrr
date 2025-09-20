// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

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

func Test_WithMeta(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		}
		e := &Error{}

		// --- When ---
		WithMeta(m)(e)

		// --- Then ---
		want := map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		}
		assert.Equal(t, want, e.meta)
	})

	t.Run("not supported types are removed", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{"A": 1, "B": struct{}{}}
		e := &Error{}

		// --- When ---
		WithMeta(m)(e)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, e.meta)
	})

	t.Run("multiple calls work like merge", func(t *testing.T) {
		// --- Given ---
		m0 := map[string]any{"A": 1, "B": 2}
		m1 := map[string]any{"B": 3}
		e := &Error{}

		// --- When ---
		WithMeta(m0)(e)
		WithMeta(m1)(e)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1, "B": 3}, e.meta)
	})
}

func Test_WithMetaFrom(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		m := TMetaAll(map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		})
		e := &Error{}

		// --- When ---
		WithMetaFrom(m)(e)

		// --- Then ---
		want := map[string]any{
			"bool":     true,
			"string":   "abc",
			"int":      2,
			"int64":    int64(2),
			"float64":  4.2,
			"time":     time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC),
			"duration": time.Second,
		}
		assert.Equal(t, want, e.meta)
	})

	t.Run("not supported types are removed", func(t *testing.T) {
		// --- Given ---
		m := TMetaAll(map[string]any{"A": 1, "B": struct{}{}})
		e := &Error{}

		// --- When ---
		WithMeta(m)(e)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, e.meta)
	})
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
		assert.Same(t, e, errors.Unwrap(err))
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

	t.Run("wrap error without options is no-op", func(t *testing.T) {
		// --- Given ---
		e := New("msg a", "a")

		// --- When ---
		err := Wrap(e)

		// --- Then ---
		assert.Same(t, e, err)
	})
}

func Test_Wrapf(t *testing.T) {
	t.Run("no %w in format", func(t *testing.T) {
		// --- When ---
		err := Wrapf("a %s", "msg")

		// --- Then ---
		assert.NotSameType(t, &Error{}, err)
		assert.Equal(t, "a msg", err.Error())
	})

	t.Run("more than one %w in format", func(t *testing.T) {
		// --- When ---
		//goland:noinspection GoPrintFunctions
		err := Wrapf("a %w: %w", io.EOF, io.ErrUnexpectedEOF)

		// --- Then ---
		assert.NotSameType(t, &Error{}, err)
		assert.ErrorEqual(t, "a EOF: unexpected EOF", err)
	})

	t.Run("wrap nil error", func(t *testing.T) {
		// --- When ---
		//goland:noinspection GoPrintFunctions
		err := Wrapf("a %w", nil)

		// --- Then ---
		assert.NotSameType(t, &Error{}, err)
		assert.ErrorEqual(t, "a %!w(<nil>)", err)
	})

	t.Run("one %w in format", func(t *testing.T) {
		// --- When ---
		//goland:noinspection GoPrintFunctions
		err := Wrapf("an %w", ErrInvJSON)

		// --- Then ---
		var xe *Error
		assert.Type(t, &xe, err)
		assert.ErrorIs(t, ErrInvJSON, xe)
		assert.Equal(t, ECInvJSON, xe.ErrorCode())
		assert.ErrorEqual(t, "an invalid JSON", err)
	})

	t.Run("one %w in format and error without error code", func(t *testing.T) {
		// --- When ---
		//goland:noinspection GoPrintFunctions
		err := Wrapf("an %w", errors.New("error"))

		// --- Then ---
		assert.NotSameType(t, &Error{}, err)
		assert.ErrorEqual(t, "an error", err)
	})
}

func Test_Error_Error(t *testing.T) {
	t.Run("xrr error", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECode")

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg", have)
	})

	t.Run("wrapped error", func(t *testing.T) {
		// --- Given ---
		e := Wrap(errors.New("msg"), WithCode("ECode"))

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg", have)
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
		assert.Nil(t, have)
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
		want := `{"error":"msg", "code":"ECode"}`
		assert.JSON(t, want, string(data))
	})

	t.Run("with metadata", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECode", Meta().Str("key", "val").Option())

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		want := `{"error":"msg", "code":"ECode", "meta": {"key": "val"}}`
		assert.JSON(t, want, string(data))
	})

	t.Run("tree", func(t *testing.T) {
		// --- Given ---
		e := TstTreeMeta()

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		want := `{
		  "code":"ECGeneric",
		  "error":"ma3\nma2\nma1",
		  "meta":{
		    "A":7,
		    "B":"b",
		    "C":"c",
		    "D":"d",
		    "E":"e",
		    "F":"f",
		    "G":"g"
		  }
		}`
		assert.JSON(t, want, string(data))
	})
}

func Test_Error_UnmarshalJSON(t *testing.T) {
	t.Run("without code and metadata", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"error": "msg"}`)
		var e *Error

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "msg", e.Error())
		assert.Equal(t, ECGeneric, e.code)
		assert.Nil(t, e.meta)
	})

	t.Run("with code and without metadata", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"error": "msg", "code":"ECode"}`)
		var e *Error

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "msg", e.Error())
		assert.Equal(t, "ECode", e.code)
		assert.Nil(t, e.meta)
	})

	t.Run("with code and metadata", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{
			"error": "msg", 
			"code":  "ECode", 
			"meta": {
				"num": 123, 
				"tim": "2022-01-18T13:57:00Z"
			}
		}`)
		var e *Error

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "msg", e.Error())
		assert.Equal(t, "ECode", e.code)
		assert.Len(t, 2, e.meta)
		assert.Equal(t, float64(123), e.meta["num"])
		assert.Equal(t, "2022-01-18T13:57:00Z", e.meta["tim"])
	})

	t.Run("error - without the error key", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"code":"code"}`)
		var e *Error

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.ErrorIs(t, ErrInvJSONError, err)
	})

	t.Run("error - invalid format", func(t *testing.T) {
		// --- Given ---
		data := []byte(`[1, 2, 3]`)
		var e *Error

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		var target *json.UnmarshalTypeError
		assert.ErrorAs(t, &target, err)
	})

	t.Run("error - invalid JSON", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{!!!}`)
		var e *Error

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		var target *json.SyntaxError
		assert.ErrorAs(t, &target, err)
	})
}

func Test_Error_Format(t *testing.T) {
	t.Run("wrapped errors", func(t *testing.T) {
		// --- Given ---
		e0 := New("msg0", "ECode0")
		e1 := Wrap(e0, WithCode("ECode1"))
		e2 := Wrap(e1, WithCode("ECode2"))

		// --- When ---
		have := fmt.Sprintf("%+v", e2)

		// --- Then ---
		assert.Equal(t, "msg0 (ECode2)", have)
	})

	t.Run("joined errors", func(t *testing.T) {
		// --- Given ---
		e0 := New("msg0", "ECode0")
		e1 := New("msg1", "ECode1")

		// --- When ---
		have := fmt.Sprintf("%+v", errors.Join(e0, e1))

		// --- Then ---
		assert.Equal(t, "msg0\nmsg1", have)
	})
}

func Test_Error_Format_tabular(t *testing.T) {
	tt := []struct {
		testN string

		msg    string
		code   string
		format string
		want   string
	}{
		{"s", "msg", "ECode", "%s", `msg`},
		{"q", "msg", "ECode", "%q", `"msg"`},
		{"v", "msg", "ECode", "%v", `msg`},
		{"+v", "msg", "ECode", "%+v", `msg (ECode)`},
		{"T", "msg", "ECode", "%T", `*xrr.Error`},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			e := New(tc.msg, tc.code)

			// --- When ---
			have := fmt.Sprintf(tc.format, e)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_Format(t *testing.T) { /* See Test_Error_Format_tabular */ }
