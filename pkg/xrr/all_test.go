// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
)

// TstErrStd is an error used in tests.
var TstErrStd = errors.New("std tst msg")

// TErrFields represents an error implementing [Fielder] interface.
type TErrFields map[string]error

func (f TErrFields) Error() string               { return "fields error" }
func (f TErrFields) ErrFields() map[string]error { return f }

// TErrMarshalJSON represents test error struct implementing [json.Marshaler]
// interface which returns 'err' error.
type TErrMarshalJSON struct{ err error }

func (tm *TErrMarshalJSON) Error() string                { return "test error" }
func (tm *TErrMarshalJSON) MarshalJSON() ([]byte, error) { return nil, tm.err }

// TstTreeCase1 returns test error tree - case 1.
//
// Shape:
//
//	   a
//	   │
//	┌──b──┐
//	c   ┌─d─┐
//	│   │   │
//	e   f   g
func TstTreeCase1() error {
	return &Error{
		msg:  "msg b",
		code: "a",
		err: &Error{
			msg:  "msg b",
			code: "b",
			err: errors.Join(
				&Error{
					msg:  "msg c",
					code: "c",
					err:  &Error{msg: "msg e", code: "e"},
				},
				&Error{
					msg:  "msg d",
					code: "d",
					err: errors.Join(
						&Error{msg: "msg f", code: "f"},
						&Error{msg: "msg g", code: "g"},
					),
				},
			),
		},
	}
}

// TstTreeCase2 returns test error tree - case 2.
//
// Shape:
//
//	   a
//	   │
//	┌──b──┐
//	c   ┌─d─┐
//	│   │   │
//	e   f   g
//	│       │
//	h       i
func TstTreeCase2() error {
	return &Error{
		msg:  "msg a",
		code: "a",
		err: &Error{
			msg:  "msg b",
			code: "b",
			err: errors.Join(
				&Error{
					msg:  "msg c",
					code: "c",
					err: &Error{
						msg:  "msg e",
						code: "e",
						err:  &Error{msg: "msg h", code: "h"},
					},
				},
				&Error{
					msg:  "msg d",
					code: "d",
					err: errors.Join(
						&Error{msg: "msg f", code: "f"},
						&Error{
							msg:  "msg g",
							code: "g",
							err:  &Error{msg: "msg i", code: "i"},
						},
					),
				},
			),
		},
	}
}

// TstTreeCase3 returns test error tree - case 3.
//
// Shape:
//
//	   a
//	   │
//	┌──b──┐
//	c   ┌─d─┐
//	│   │   │
//	e   f   g
//	│   │
//	h   i
func TstTreeCase3() error {
	return &Error{
		msg:  "msg a",
		code: "a",
		err: &Error{
			msg:  "msg b",
			code: "b",
			err: errors.Join(
				&Error{
					msg:  "msg c",
					code: "c",
					err: &Error{
						msg:  "msg e",
						code: "e",
						err:  &Error{msg: "msg h", code: "h"},
					},
				},
				&Error{
					msg:  "msg d",
					code: "d",
					err: errors.Join(
						&Error{
							msg:  "msg f",
							code: "f",
							err:  &Error{msg: "msg f", code: "i"},
						},
						&Error{msg: "msg g", code: "g"},
					),
				},
			),
		},
	}
}

// TstTreeCase4 returns test error tree - case 4.
//
// Shape:
//
//	┌─────┐
//	a   ┌─b─┐
//	│   │   │
//	c   d   e
func TstTreeCase4() error {
	return errors.Join(
		&Error{
			msg:  "msg a",
			code: "a",
			err:  &Error{msg: "msg c", code: "c"},
		},
		&Error{
			msg:  "msg b",
			code: "b",
			err: errors.Join(
				&Error{msg: "msg d", code: "d"},
				&Error{msg: "msg e", code: "e"},
			),
		},
	)
}
