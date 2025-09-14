// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

// AssertError asserts that the provided error is non-nil and is an instance of
// [xrr.Error]. If the error is an instance of [xrr.Error], it returns true and
// the [xrr.Error] instance. If the error is nil or not an instance of
// [xrr.Error], it marks the test as failed, logs an error message, and returns
// false and nil.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Error]
// without unwrapping.
func AssertError(t tester.T, err error) (*xrr.Error, bool) {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return nil, false // nolint: nilerr
	}
	var xe *xrr.Error
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

// AssertEqual asserts that the provided error is non-nil and is an instance of
// [xrr.Error], then asserts the error has the wanted error message and error
// code. Returns true if it is, otherwise marks the test as failed, writes an
// error message to the test log, and returns false.
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
// it does otherwise, marks the test as failed, writes an error message to the
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
// Returns true if the key doesn't exist otherwise, marks the test as failed,
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
// true if it does otherwise, marks the test as failed, writes an error message
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
// true if it has otherwise, marks the test as failed, writes an error message
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
// true if it has otherwise, marks the test as failed, writes an error message
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
// true if it has otherwise, marks the test as failed, writes an error message
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

// AssertFields asserts that the provided error is non-nil and is an instance
// of [xrr.Fields]. If the error is an instance of [xrr.Fields], it returns
// true and the [xrr.Fields] instance. If the error is nil or not an instance
// of [xrr.Error], it marks the test as failed, logs an error message, and
// returns false and nil.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
// without unwrapping.
func AssertFields(t tester.T, err error) (xrr.Fields, bool) {
	t.Helper()
	if e := check.NotNil(err); e != nil {
		t.Error(notice.From(e).SetHeader("[xrr] expected error not to be nil"))
		return nil, false // nolint: nilerr
	}
	var xe xrr.Fields
	if e := check.Type(&xe, err); e != nil {
		msg := notice.From(e).
			SetHeader("[xrr] expected xrr.Fields instance").
			Remove("src").
			Append("error", "%T", err)
		t.Error(msg)
		return nil, false // nolint: nilerr
	}
	return xe, true
}

// AssertFieldsEqual asserts that the provided error is non-nil and is an
// instance of [xrr.Fields]. Then asserts the string error message equals to
// the one provided. Returns true if it does, otherwise marks the test as
// failed, writes an error message to the test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
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

// AssertFieldCnt asserts that the provided error is non-nil and is an instance
// of [xrr.Fields]. Then asserts it has the given number of fields. Returns
// true if it has, otherwise marks the test as failed, writes an error message
// to the test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
// without unwrapping.
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

// AssertHasField asserts error is an instance of [xrr.Fields] and has the
// given field name. Returns the value of the error (might be nil) and true
// if the field exists. Otherwise, marks the test as failed, writes an error
// message to the test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
// without unwrapping.
func AssertHasField(t tester.T, field string, err error) (error, bool) {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return nil, false
	}
	ve, e := check.HasKey(field, xe)
	if e != nil {
		msg := notice.From(e, "xrr").
			Remove("map").
			Append("fields", "%s", dump.New().Any(xe))
		t.Error(msg)
		return nil, false
	}
	return ve, true
}

// AssertFieldEqual asserts error is an instance of [xrr.Fields] and has the
// given field name and the error message equals to msg. Returns true if it
// does, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
// without unwrapping.
func AssertFieldEqual(t tester.T, field, want string, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}
	ve, e := check.HasKey(field, xe)
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

// AssertFieldCode asserts error is an instance of [xrr.Fields] and has the
// given field name with an error having the given error code. Returns true if
// it does, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
// without unwrapping.
func AssertFieldCode(t tester.T, field, code string, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}

	ve, e := check.HasKey(field, xe)
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

// AssertFieldIs asserts error is an instance of [xrr.Fields] and has the
// given field name with an error which has "want" error in its chain.
// Assertion uses [errors.Is] to check the field error chain. Returns true if
// it does, otherwise marks the test as failed, writes an error message to the
// test log, and returns false.
//
// Unlike [errors.As], it directly checks if the error is of type [xrr.Fields]
// without unwrapping.
func AssertFieldIs(t tester.T, field string, want, err error) bool {
	t.Helper()
	xe, success := AssertFields(t, err)
	if !success {
		return false
	}
	ve, e := check.HasKey(field, xe)
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
