// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_ErrorFactory(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		// --- Given ---
		have := ErrorFunc[EDXrr]()

		// --- When ---
		err := have("msg", "ECode")

		// --- Then ---
		e, _ := assert.SameType(t, &GenericError[EDXrr]{}, err)
		assert.Equal(t, "msg", e.msg)
		assert.Equal(t, "ECode", e.code)
		assert.Nil(t, e.meta)
		assert.Nil(t, e.err)
	})

	t.Run("with options", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{"A": 1, "B": func() {}}
		have := ErrorFunc[EDXrr]()

		// --- When ---
		err := have("msg", "ECode", WithMeta(m))

		// --- Then ---
		e, _ := assert.SameType(t, &GenericError[EDXrr]{}, err)
		assert.Equal(t, "msg", e.msg)
		assert.Equal(t, "ECode", e.code)
		assert.NotSame(t, m, e.meta)
		assert.Equal(t, map[string]any{"A": 1}, e.meta)
		assert.Nil(t, e.err)
	})
}

func Test_GenericError_Error(t *testing.T) {
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

	t.Run("case 1", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase1()

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg e; msg f; msg g", have)
	})

	t.Run("case 2", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase2()

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg h; msg f; msg i", have)
	})

	t.Run("case 3", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg h; msg f; msg g", have)
	})

	t.Run("case 4", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase4()

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg c\nmsg d; msg e", have)
	})

	t.Run("case 5", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase5()

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "a: msg b; d: msg e; f: msg g; msg h", have)
	})

	t.Run("instance with wrapped error", func(t *testing.T) {
		// --- Given ---
		errFunc := ErrorFunc[string]()
		e := errFunc("msg", "ECode", WithCause(errors.New("cause")))

		// --- When ---
		have := e.Error()

		// --- Then ---
		assert.Equal(t, "msg: cause", have)
	})
}

func Test_GenericError_ErrorCode(t *testing.T) {
	t.Run("returns code", func(t *testing.T) {
		// --- Given ---
		e := &GenericError[string]{code: "ECode"}

		// --- When ---
		have := e.ErrorCode()

		// --- Then ---
		assert.Equal(t, "ECode", have)
	})

	t.Run("returns ECGeneric when code is empty", func(t *testing.T) {
		// --- Given ---
		e := &GenericError[string]{}

		// --- When ---
		have := e.ErrorCode()

		// --- Then ---
		assert.Equal(t, ECGeneric, have)
	})
}

func Test_GenericError_MetaAll(t *testing.T) {
	t.Run("returns a clone", func(t *testing.T) {
		// --- Given ---
		m := map[string]any{"A": 1, "B": 2}
		e := &GenericError[string]{meta: m}

		// --- When ---
		have := e.MetaAll()

		// --- Then ---
		assert.NotSame(t, m, have)
		assert.Equal(t, m, have)
	})
}

func Test_GenericError_Unwrap(t *testing.T) {
	t.Run("returns wrapped error", func(t *testing.T) {
		// --- Given ---
		e := New("msg", "ECode")

		// --- When ---
		have := errors.Unwrap(e)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("returns nil for nil instance", func(t *testing.T) {
		// --- Given ---
		var e *GenericError[string]

		// --- When ---
		have := e.Unwrap()

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_GenericError_MarshalJSON(t *testing.T) {
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
		  "error":"ma3; ma2; ma1",
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

	t.Run("metadata keys are alphabetically sorted", func(t *testing.T) {
		// --- Given ---
		meta := Meta().Str("zebra", "z").Str("apple", "a").Str("mango", "m")
		e := New("msg", "ECode", meta.Option())

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		want := `{
			"code": "ECode",
			"error": "msg",
			"meta": {
				"apple": "a",
				"mango": "m",
				"zebra": "z"
			}
		}`
		assert.JSON(t, want, string(data))
	})
}

func Test_GenericError_UnmarshalJSON(t *testing.T) {
	t.Run("without code and metadata", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"error": "msg"}`)
		var e *GenericError[string]

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
		var e *GenericError[string]

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
		var e *GenericError[string]

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
		var e *GenericError[string]

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.ErrorIs(t, ErrInvJSONError, err)
	})

	t.Run("error - invalid format", func(t *testing.T) {
		// --- Given ---
		data := []byte(`[1, 2, 3]`)
		var e *GenericError[string]

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		var target *json.UnmarshalTypeError
		assert.ErrorAs(t, &target, err)
	})

	t.Run("error - invalid JSON", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{!!!}`)
		var e *GenericError[string]

		// --- When ---
		err := json.Unmarshal(data, &e)

		// --- Then ---
		var target *json.SyntaxError
		assert.ErrorAs(t, &target, err)
	})
}

func Test_GenericError_Format(t *testing.T) {
	t.Run("wrapped errors", func(t *testing.T) {
		// --- Given ---
		e0 := New("msg0", "ECode0")
		e1 := Wrap(e0, WithCode("ECode1"))
		e2 := WrapUsing[string](e1, WithCode("ECode2"))

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

func Test_GenericError_Format_tabular(t *testing.T) {
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
		{
			"T",
			"msg",
			"ECode",
			"%T",
			`*xrr.GenericError[github.com/ctx42/xrr/pkg/xrr.EDXrr]`,
		},
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

func Test_Format(t *testing.T) { /* See Test_GenericError_Format_tabular */ }
