// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrrtest

import (
	"fmt"
	"time"

	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"

	"github.com/ctx42/xrr/pkg/xrr"
)

// AssertError asserts err is non-nil and is an instance of
// [xrr.GenericError[T]]. Returns the [xrr.GenericError[T]] instance and true
// on success. If err is nil or not an instance of [xrr.GenericError[T]], it
// marks the test as failed, writes an error message to the test log, and
// returns nil and false.
//
// Unlike [errors.As], it directly checks if the error is of the type
// [xrr.GenericError[T]] without unwrapping.
func AssertError[T xrr.Domain](t tester.T, err error) (*xrr.GenericError[T], bool) {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return nil, false // nolint: nilerr
	}
	// Target variable for the type assertion into *xrr.GenericError[T].
	var xe *xrr.GenericError[T]
	if e := check.Type(&xe, err); e != nil {
		msg := notice.From(e).
			SetHeader("[xrr] expected *xrr.Error instance").
			Remove("src").
			Append("error", "%T", err)
		t.Error(msg)
		return nil, false // nolint: nilerr
	}
	return xe, true
}

// AssertEqual asserts err is non-nil and its %+v formatted representation
// equals "want". For [xrr.GenericError], %+v includes the error message
// followed by the error code in parentheses. Returns true on success,
// otherwise marks the test as failed, writes an error message to the test log,
// and returns false.
func AssertEqual(t tester.T, want string, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have := fmt.Sprintf("%+v", err)
	if e := check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error to have a message"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertCode asserts err is not nil, implements [xrr.Coder] interface, and has
// the given error code. Returns true if it does, otherwise marks the test as
// failed, writes an error message to the test log, and returns false.
func AssertCode(t tester.T, want string, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	var coder xrr.Coder
	if e := check.Type(&coder, err); e != nil {
		const hHeader = "[xrr] expected xrr.Coder instance"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e := check.Equal(want, coder.ErrorCode()); e != nil {
		const hHeader = "[xrr] expected error with error code"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertKeyCnt asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has a given number of keys. Returns true if
// it does, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
func AssertKeyCnt(t tester.T, want int, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	if e := check.Len(want, xrr.GetMeta(err)); e != nil {
		const hHeader = "[xrr] expected error number of metadata keys"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertNoKey asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], does not contain the key with the given name.
// Returns true if the key doesn't exist, otherwise marks the test as failed,
// writes an error message to the test log, and returns false.
func AssertNoKey(t tester.T, key string, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	if e := check.HasNoKey(key, xrr.GetMeta(err)); e != nil {
		const hHeader = "[xrr] expected error without the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertStr asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has the key with the given value. Returns
// true if it does, otherwise marks the test as failed, writes an error message
// to the test log, and returns false.
func AssertStr(t tester.T, key, want string, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have, e := check.HasKey(key, xrr.GetMeta(err))
	if e != nil {
		const hHeader = "[xrr] expected error to have the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e = check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertInt asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has the key with the given value. Returns
// true if it has, otherwise marks the test as failed, writes an error
// message to the test log, and returns false.
func AssertInt(t tester.T, key string, want int, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have, e := check.HasKey(key, xrr.GetMeta(err))
	if e != nil {
		const hHeader = "[xrr] expected error to have the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e = check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertInt64 asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has the key with the given value. Returns
// true if it has, otherwise marks the test as failed, writes an error message
// to the test log, and returns false.
func AssertInt64(t tester.T, key string, want int64, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have, e := check.HasKey(key, xrr.GetMeta(err))
	if e != nil {
		const hHeader = "[xrr] expected error to have the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e = check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertFloat64 asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has the key with the given value. Returns
// true if it has, otherwise marks the test as failed, writes an error message
// to the test log, and returns false.
func AssertFloat64(t tester.T, key string, want float64, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have, e := check.HasKey(key, xrr.GetMeta(err))
	if e != nil {
		const hHeader = "[xrr] expected error to have the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e = check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertBool asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has the key with the given value. Returns
// true if it has, otherwise marks the test as failed, writes an error message
// to the test log, and returns false.
func AssertBool(t tester.T, key string, want bool, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have, e := check.HasKey(key, xrr.GetMeta(err))
	if e != nil {
		const hHeader = "[xrr] expected error to have the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e = check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertTime asserts that the provided error is non-nil and error metadata,
// retrieved using [xrr.GetMeta], has the key with the given value. Returns
// true if it has, otherwise marks the test as failed, writes an error message
// to the test log, and returns false.
func AssertTime(t tester.T, key string, want time.Time, err error) bool {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return false
	}
	have, e := check.HasKey(key, xrr.GetMeta(err))
	if e != nil {
		const hHeader = "[xrr] expected error to have the metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	if e = check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error metadata key"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertFields asserts err is non-nil and implements [xrr.Fielder]. Returns
// the [xrr.Fielder] instance and true on success. If err is nil or does not
// implement [xrr.Fielder], it marks the test as failed, writes an error
// message to the test log, and returns nil and false.
//
// Unlike [errors.As], it directly checks if the error implements [xrr.Fielder]
// without unwrapping.
func AssertFields(t tester.T, err error) (xrr.Fielder, bool) {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return nil, false // nolint: nilerr
	}
	// Target variable for the type assertion into xrr.Fielder.
	var xe xrr.Fielder
	if e := check.Type(&xe, err); e != nil {
		msg := notice.From(e).
			SetHeader("[xrr] expected xrr.Fielder instance").
			Remove("src").
			Append("error", "%T", err)
		t.Error(msg)
		return nil, false // nolint: nilerr
	}
	return xe, true
}

// AssertFieldsEqual asserts err is non-nil and implements [xrr.Fielder], then
// asserts its %+v formatted representation equals "want". Returns true on
// success, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error implements [xrr.Fielder]
// without unwrapping.
func AssertFieldsEqual(t tester.T, want string, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}
	have := fmt.Sprintf("%+v", xe)
	if e := check.Equal(want, have); e != nil {
		const hHeader = "[xrr] expected error to have a message"
		t.Error(notice.From(e).SetHeader(hHeader))
		return false
	}
	return true
}

// AssertFieldCnt asserts err is non-nil and implements [xrr.Fielder], then
// asserts it has the given number of fields. Returns true on success,
// otherwise marks the test as failed, writes an error message to the test
// log, and returns false.
//
// Unlike [errors.As], it directly checks if the error implements
// [xrr.Fielder] without unwrapping.
func AssertFieldCnt(t tester.T, want int, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}
	if e := check.Len(want, xe); e != nil {
		msg := notice.From(e, "xrr").Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	return true
}

// AssertHasField asserts error is an instance of [xrr.Fielder] and has the
// given field name. Returns the field's error value (might be nil) and true
// on success. Otherwise, marks the test as failed, writes an error message
// to the test log, and returns nil and false.
//
// Unlike [errors.As], it directly checks if the error is of type
// [xrr.Fielder] without unwrapping.
func AssertHasField(t tester.T, field string, err error) (error, bool) {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return nil, false
	}
	// Verify the field exists and retrieve its error value.
	ve, e := check.HasKey(field, xe.ErrorFields())
	if e != nil {
		msg := notice.From(e, "xrr").
			Remove("map").
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return nil, false
	}
	return ve, true
}

// AssertFieldEqual asserts error is an instance of [xrr.Fielder] and has the
// given field name with an error message that equals "want". Returns true on
// success, otherwise marks the test as failed, writes an error message to
// the test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type
// [xrr.Fielder] without unwrapping.
func AssertFieldEqual(t tester.T, field, want string, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}
	// Verify the field exists and retrieve its error value.
	ve, e := check.HasKey(field, xe.ErrorFields())
	if e != nil {
		msg := notice.From(e, "xrr").
			Remove("map").
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	if e = check.ErrorEqual(want, ve); e != nil {
		msg := notice.From(e, "xrr").
			Remove("map").
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	return true
}

// AssertFieldCode asserts error is an instance of [xrr.Fielder] and has the
// given field name with an error having the given error code. Returns true if
// it does, otherwise marks the test as failed, writes an error message to
// the test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fielder]
// without unwrapping.
func AssertFieldCode(t tester.T, field, code string, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}

	// Verify the field exists and retrieve its error value.
	ve, e := check.HasKey(field, xe.ErrorFields())
	if e != nil {
		hHeader := "[xrr] expected field to exist"
		msg := notice.From(e).
			SetHeader(hHeader).
			Remove("key").
			Prepend("field", "%s", field).
			Remove("map").
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	if e = check.Equal(code, xrr.GetCode(ve)); e != nil {
		const hHeader = "[xrr] expected field to have the given error code"
		msg := notice.From(e).
			SetHeader(hHeader).
			Prepend("field", "%s", field).
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	return true
}

// AssertFieldIs asserts err is an instance of [xrr.Fielder] and that the field
// identified by field exists with an error that has "want" in its chain. The
// assertion uses [errors.Is] to check the field error chain. Returns true on
// success, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type
// [xrr.Fielder] without unwrapping.
func AssertFieldIs(t tester.T, field string, want, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}
	// Verify the field exists and retrieve its error value.
	ve, e := check.HasKey(field, xe.ErrorFields())
	if e != nil {
		msg := notice.From(e, "xrr").
			Remove("map").
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	if e = check.ErrorIs(want, ve); e != nil {
		msg := notice.From(e, "xrr").Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return false
	}
	return true
}
