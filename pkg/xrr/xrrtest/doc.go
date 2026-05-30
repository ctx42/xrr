// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package xrrtest provides testing helpers for the xrr error library.
//
// It is built on top of github.com/ctx42/testing and offers precise,
// domain-aware assertions for [xrr.GenericError], [xrr.GenericFields],
// metadata values, and field validation errors.
//
// Most helpers use direct type assertions (not errors.As) so that test
// failures clearly show the expected vs actual values.
//
// Primary groups of helpers:
//   - AssertError, AssertEqual, AssertCode, AssertKeyCnt, AssertNoKey
//   - AssertStr, AssertInt, AssertInt64, AssertFloat64, AssertBool,
//     AssertTime, AssertDuration (metadata assertions)
//   - AssertFields, AssertFieldsEqual, AssertFieldCnt, AssertHasField,
//     AssertFieldEqual, AssertFieldCode, AssertFieldIs (field error helpers)
//
// See assertions.go for the complete set and usage examples.
package xrrtest
