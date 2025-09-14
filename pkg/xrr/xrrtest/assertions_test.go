// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrrtest

import (
	"errors"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/tester"

	"github.com/ctx42/xrr/pkg/xrr"
)

func Test_AssertError(t *testing.T) {
	t.Run("success - error is instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		want := xrr.New("msg", "ECode")

		// --- When ---
		have, success := AssertError(tspy, want)

		// --- Then ---
		assert.True(t, success)
		assert.Same(t, want, have)
	})

	t.Run("error - error not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have, success := AssertError(tspy, err)

		// --- Then ---
		assert.False(t, success)
		assert.Nil(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have, success := AssertError(tspy, nil)

		// --- Then ---
		assert.False(t, success)
		assert.Nil(t, have)
	})
}

func Test_AssertMsg(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertMsg(tspy, "msg", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertMsg(tspy, "key", nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - message not equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have the message:\n" +
			"  want: \"other\"\n" +
			"  have: \"msg\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertMsg(tspy, "other", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertEqual(tspy, "msg (ECode)", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertEqual(tspy, "key", nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - message and code not equal", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have a message:\n" +
			"  want: \"other (ECode)\"\n" +
			"  have: \"msg (ECode)\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertEqual(tspy, "other (ECode)", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertCode(tspy, "ECode", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertCode(tspy, "key", nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Coder", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Coder instance:\n" +
			"  target: xrr.Coder\n" +
			"     src: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertCode(tspy, "ECode", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - code does not match", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error with error code:\n" +
			"  want: \"ECOther\"\n" +
			"  have: \"ECode\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertCode(tspy, "ECOther", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertKeys(t *testing.T) {
	t.Run("no metadata keys", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertKeys(tspy, 0, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("some metadata keys", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		meta := xrr.Meta().Str("str", "a").Int("int", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertKeys(tspy, 2, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - different number of metadata keys", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		wMsg := "" +
			"[xrr] expected error number of metadata keys:\n" +
			"  want: 3\n" +
			"  have: 2"
		tspy.ExpectLogEqual(wMsg)
		tspy.ExpectError()
		tspy.Close()

		meta := xrr.Meta().Str("str", "a").Int("int", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertKeys(tspy, 3, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertKeys(tspy, 1, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertNoKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertNoKey(tspy, "key", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertNoKey(tspy, "key", nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error has a metadata key", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error without the metadata key:\n" +
			"    key: \"key\"\n" +
			"  value: \"val\"\n" +
			"    map:\n" +
			"         map[string]any{\n" +
			"           \"key\": \"val\",\n" +
			"         }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Str("key", "val")
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertNoKey(tspy, "key", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertNoKey(tspy, "key", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertStr(t *testing.T) {
	t.Run("success - error has the key value pair", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		meta := xrr.Meta().Str("key", "val")
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertStr(tspy, "key", "val", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertStr(tspy, "key", "val", nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have metadata key:\n" +
			"  key: \"key\"\n" +
			"  map:\n" +
			"       map[string]any{\n" +
			"         \"other\": \"val\",\n       }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Str("other", "val")
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertStr(tspy, "key", "val", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key is not of the string type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want type: string\n" +
			"  have type: int"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertStr(tspy, "key", "val", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertStr(tspy, "key", "val", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertInt(t *testing.T) {
	t.Run("success - error has the key value pair", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		meta := xrr.Meta().Int("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertInt(tspy, "key", 1, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertInt(tspy, "key", 1, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have metadata key:\n" +
			"  key: \"key\"\n" +
			"  map:\n" +
			"       map[string]any{\n" +
			"         \"other\": 1,\n" +
			"       }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int("other", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertInt(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key is not of the int type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want type: int\n" +
			"  have type: int64"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int64("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertInt(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertInt(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertInt64(t *testing.T) {
	t.Run("success - error has the key value pair", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		meta := xrr.Meta().Int64("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertInt64(tspy, "key", 1, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertInt64(tspy, "key", 1, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have metadata key:\n" +
			"  key: \"key\"\n" +
			"  map:\n" +
			"       map[string]any{\n" +
			"         \"other\": 1,\n" +
			"       }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int64("other", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertInt64(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key is not of the int64 type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want type: int64\n" +
			"  have type: int"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertInt64(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertInt64(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertFloat64(t *testing.T) {
	t.Run("success - error has the key value pair", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		meta := xrr.Meta().Float64("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertFloat64(tspy, "key", 1, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertFloat64(tspy, "key", 1, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have metadata key:\n" +
			"  key: \"key\"\n" +
			"  map:\n" +
			"       map[string]any{\n" +
			"         \"other\": 1,\n" +
			"       }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Float64("other", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertFloat64(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key is not of the float64 type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want type: float64\n" +
			"  have type: int"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertFloat64(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertFloat64(tspy, "key", 1, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertBool(t *testing.T) {
	t.Run("success - error has the key value pair", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		meta := xrr.Meta().Bool("key", true)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertBool(tspy, "key", true, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertBool(tspy, "key", true, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have metadata key:\n" +
			"  key: \"key\"\n" +
			"  map:\n" +
			"       map[string]any{\n" +
			"         \"other\": true,\n" +
			"       }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Bool("other", true)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertBool(tspy, "key", true, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key is not of the bool type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want type: bool\n" +
			"  have type: int"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertBool(tspy, "key", true, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error is not an instance of xrr.Error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected *xrr.Error instance:\n" +
			"  target: *xrr.Error\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertBool(tspy, "key", true, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertTime(t *testing.T) {
	t.Run("success - error has the key value pair", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		tim := time.Date(2022, 3, 18, 0, 0, 0, 0, time.UTC)
		meta := xrr.Meta().Time("key", tim)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		have := AssertTime(tspy, "key", tim, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertTime(tspy, "key", time.Time{}, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have metadata key:\n" +
			"  key: \"key\"\n" +
			"  map: map[string]any(nil)"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.New("msg", "ECode")

		// --- When ---
		have := AssertTime(tspy, "key", time.Now(), err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - timezone does not match", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want: 2022-03-18T00:00:00Z\n" +
			"  have: 2022-03-17T23:00:00Z ( 2022-03-18T00:00:00+01:00 )\n" +
			"  diff: 1h0m0s"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		tz, _ := time.LoadLocation("Europe/Warsaw")
		tim := time.Date(2022, 3, 18, 0, 0, 0, 0, tz)
		meta := xrr.Meta().Time("key", tim)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		exp := time.Date(2022, 3, 18, 0, 0, 0, 0, time.UTC)
		have := AssertTime(tspy, "key", exp, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - key is not of the time type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error metadata key:\n" +
			"  want type: time.Time\n" +
			"  have type: int"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		meta := xrr.Meta().Int("key", 1)
		err := xrr.New("msg", "ECode", meta.Option())

		// --- When ---
		exp := time.Date(2022, 3, 18, 0, 0, 0, 0, time.UTC)
		have := AssertTime(tspy, "key", exp, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertFields(t *testing.T) {
	t.Run("success - error is instance of xrr.Fields", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.Fields{"a": errors.New("msg")}

		// --- When ---
		have, success := AssertFields(tspy, err)

		// --- Then ---
		assert.True(t, success)
		assert.Same(t, err, have)
	})

	t.Run("error - error not an instance of xrr.Fields", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Fields instance:\n" +
			"  target: xrr.Fields\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have, success := AssertFields(tspy, err)

		// --- Then ---
		assert.False(t, success)
		assert.Nil(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have, success := AssertFields(tspy, nil)

		// --- Then ---
		assert.False(t, success)
		assert.Nil(t, have)
	})
}

func Test_AssertFieldCnt(t *testing.T) {
	t.Run("success - number of fields does match", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": errors.New("m1"),
		}

		// --- When ---
		have := AssertFieldCnt(tspy, 2, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertFieldCnt(tspy, 2, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - number of fields does not match", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Fields length:\n" +
			"    want: 3\n" +
			"    have: 2\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f0\": \"m0\",\n" +
			"            \"f1\": \"m1\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": errors.New("m1"),
		}

		// --- When ---
		have := AssertFieldCnt(tspy, 3, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error not an instance of xrr.Fields", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Fields instance:\n" +
			"  target: xrr.Fields\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertFieldCnt(tspy, 1, err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertHasField(t *testing.T) {
	t.Run("success - field exists", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		f0 := errors.New("m0")
		f1 := errors.New("m1")
		err := xrr.Fields{"f0": f0, "f1": f1}

		// --- When ---
		have, success := AssertHasField(tspy, "f0", err)

		// --- Then ---
		assert.Same(t, f0, have)
		assert.True(t, success)
	})

	t.Run("success - field exists and it is nil", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		f0 := errors.New("m0")
		err := xrr.Fields{"f0": f0, "f1": nil}

		// --- When ---
		have, success := AssertHasField(tspy, "f1", err)

		// --- Then ---
		assert.Nil(t, have)
		assert.True(t, success)
	})

	t.Run("error - field does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected map to have a key:\n" +
			"     key: \"f1\"\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f0\": \"m0\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.Fields{"f0": errors.New("m0")}

		// --- When ---
		have, success := AssertHasField(tspy, "f1", err)

		// --- Then ---
		assert.Nil(t, have)
		assert.False(t, success)
	})

	t.Run("error - error not an instance of xrr.Fields", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Fields instance:\n" +
			"  target: xrr.Fields\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have, success := AssertHasField(tspy, "f1", err)

		// --- Then ---
		assert.False(t, success)
		assert.Nil(t, have)
	})
}

func Test_AssertFieldEqual(t *testing.T) {
	t.Run("susses - field exists and with given message", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": xrr.New("m1", "ECF1"),
		}

		// --- When ---
		have := AssertFieldEqual(tspy, "f1", "m1", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		have := AssertFieldEqual(tspy, "f1", "m1", nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - field does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected map to have a key:\n" +
			"  key: \"f2\"\n" +
			"  map:\n" +
			"       map[string]error{\n" +
			"         \"f0\": \"m0\",\n" +
			"         \"f1\": \"m1\",\n" +
			"       }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": xrr.New("m1", "ECF1"),
		}

		// --- When ---
		have := AssertFieldEqual(tspy, "f2", "m2", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - filed exists but message is different", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected the error message to be:\n" +
			"    want: \"other\"\n" +
			"    have: \"m1\"\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f0\": \"m0\",\n" +
			"            \"f1\": \"m1\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": xrr.New("m1", "ECF1"),
		}

		// --- When ---
		have := AssertFieldEqual(tspy, "f1", "other", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - error not an instance of xrr.Fields", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Fields instance:\n" +
			"  target: xrr.Fields\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertFieldEqual(tspy, "f1", "msg", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertFieldCode(t *testing.T) {
	t.Run("success - field exists with given error code", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": xrr.New("m1", "ECF1"),
		}

		// --- When ---
		have := AssertFieldCode(tspy, "f1", "ECF1", err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - field exists with different error code", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected field to have the given error code:\n" +
			"   field: f1\n" +
			"    want: \"other\"\n" +
			"    have: \"ECF1\"\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f0\": \"m0\",\n" +
			"            \"f1\": \"m1\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
			"f1": xrr.New("m1", "ECF1"),
		}

		// --- When ---
		have := AssertFieldCode(tspy, "f1", "other", err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - field does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected field to exist:\n" +
			"   field: f2\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f0\": \"m0\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := xrr.Fields{
			"f0": errors.New("m0"),
		}

		// --- When ---
		have := AssertFieldCode(tspy, "f2", "other", err)

		// --- Then ---
		assert.False(t, have)
	})
	t.Run("error - error not an instance of xrr.Fields", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected xrr.Fields instance:\n" +
			"  target: xrr.Fields\n" +
			"   error: *errors.errorString"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		err := errors.New("some error")

		// --- When ---
		have := AssertFieldCode(tspy, "f1", "code", err)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_AssertFieldIs(t *testing.T) {
	t.Run("success - field exists and has error in chain", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		want := xrr.New("m1", "ECF1")
		err := xrr.Fields{
			"f1": want,
		}

		// --- When ---
		have := AssertFieldIs(tspy, "f1", want, err)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("error - nil error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "[xrr] expected error not to be nil"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		want := xrr.New("m1", "ECF1")
		have := AssertFieldIs(tspy, "f1", want, nil)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - field does not exist", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected map to have a key:\n" +
			"     key: \"other\"\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f0\": \"m1\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		want := xrr.New("m1", "ECF1")
		err := xrr.Fields{
			"f0": want,
		}

		// --- When ---
		have := AssertFieldIs(tspy, "other", want, err)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("error - field exists but error not in chain", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		wMsg := "" +
			"[xrr] expected error to have a target in its tree:\n" +
			"    want: (*xrr.Error) m1\n" +
			"    have: (*xrr.Error) m1\n" +
			"  fields:\n" +
			"          map[string]error{\n" +
			"            \"f1\": \"m1\",\n" +
			"          }"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		want := xrr.New("m1", "ECF1")
		err := xrr.Fields{
			"f1": xrr.New("m1", "ECF1"),
		}

		// --- When ---
		have := AssertFieldIs(tspy, "f1", want, err)

		// --- Then ---
		assert.False(t, have)
	})
}
