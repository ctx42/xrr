// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
)

// ErrTst is an error used in tests.
var ErrTst = errors.New("std tst msg")

// TErrorFields represents an error implementing [Fielder] interface.
type TErrorFields map[string]error

func (f TErrorFields) Error() string                 { return "fields error" }
func (f TErrorFields) ErrorFields() map[string]error { return f }

// TFielderCoder represents an error implementing both [Fielder] and [Coder].
type TFielderCoder struct {
	code   string
	fields map[string]error
}

func (t TFielderCoder) Error() string                 { return "fielder coder" }
func (t TFielderCoder) ErrorCode() string             { return t.code }
func (t TFielderCoder) ErrorFields() map[string]error { return t.fields }

// TErrMarshalJSON represents a test error struct implementing [json.Marshaler]
// interface which returns 'err' error.
type TErrMarshalJSON struct{ err error }

func (tm *TErrMarshalJSON) Error() string                { return "test error" }
func (tm *TErrMarshalJSON) MarshalJSON() ([]byte, error) { return nil, tm.err }

// TMetaAll represents a struct implementing [Metadater] interface.
type TMetaAll map[string]any

func (T TMetaAll) MetaAll() map[string]any { return T }

// TstMetaTree returns a test error tree.
//
// Shape (metadata):
//
//	     A7,B0
//	       │
//	  ┌──A6,C1───┐
//	  │          │
//	A4,D2    ┌─A5,E2─┐
//	  │      │       │
//	A3,F3  A2,G3   A1,H3
func TstMetaTree() error {
	a3f3 := New("A3,F3", "A3", Meta().Int("A", 3).Int("F", 3).Option())
	a4d2 := Wrap(a3f3, Meta().Int("A", 4).Int("D", 2).Option())

	a1h3 := New("A1,H3", "A1", Meta().Int("A", 1).Int("H", 3).Option())
	a2g3 := New("A2,H3", "A2", Meta().Int("A", 2).Int("G", 3).Option())

	a5e2 := errors.Join(a1h3, a2g3)
	a5e2 = Wrap(a5e2, Meta().Int("A", 5).Int("E", 2).Option())

	a6c1 := errors.Join(a4d2, a5e2)
	a6c1 = Wrap(a6c1, Meta().Int("A", 6).Int("C", 1).Option())

	a7b0 := Wrap(a6c1, Meta().Int("A", 7).Int("B", 0).Option())

	return a7b0
}

// TstTreeCase1 returns a test error tree - case 1.
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
	return &GenericError[EDXrr]{
		code: "a",
		err: &GenericError[EDXrr]{
			code: "b",
			err: errors.Join(
				&GenericError[EDXrr]{
					code: "c",
					err:  &GenericError[EDXrr]{msg: "msg e", code: "e"},
				},
				&GenericError[EDXrr]{
					code: "d",
					err: errors.Join(
						&GenericError[EDXrr]{msg: "msg f", code: "f"},
						&GenericError[EDXrr]{msg: "msg g", code: "g"},
					),
				},
			),
		},
	}
}

// TstTreeCase2 returns a test error tree - case 2.
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
	return &GenericError[EDXrr]{
		code: "a",
		err: &GenericError[EDXrr]{
			code: "b",
			err: errors.Join(
				&GenericError[EDXrr]{
					code: "c",
					err: &GenericError[EDXrr]{
						code: "e",
						err:  &GenericError[EDXrr]{msg: "msg h", code: "h"},
					},
				},
				&GenericError[EDXrr]{
					code: "d",
					err: errors.Join(
						&GenericError[EDXrr]{msg: "msg f", code: "f"},
						&GenericError[EDXrr]{
							code: "g",
							err:  &GenericError[EDXrr]{msg: "msg i", code: "i"},
						},
					),
				},
			),
		},
	}
}

// TstTreeCase3 returns a test error tree - case 3.
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
	return &GenericError[EDXrr]{
		code: "a",
		err: &GenericError[EDXrr]{
			code: "b",
			err: errors.Join(
				&GenericError[EDXrr]{
					code: "c",
					err: &GenericError[EDXrr]{
						code: "e",
						err:  &GenericError[EDXrr]{msg: "msg h", code: "h"},
					},
				},
				&GenericError[EDXrr]{
					code: "d",
					err: errors.Join(
						&GenericError[EDXrr]{
							code: "f",
							err:  &GenericError[EDXrr]{msg: "msg f", code: "i"},
						},
						&GenericError[EDXrr]{msg: "msg g", code: "g"},
					),
				},
			),
		},
	}
}

// TstTreeCase4 returns a test error tree - case 4.
//
// Shape:
//
//	┌─────┐
//	a   ┌─b─┐
//	│   │   │
//	c   d   e
func TstTreeCase4() error {
	return errors.Join(
		&GenericError[EDXrr]{
			code: "a",
			err:  &GenericError[EDXrr]{msg: "msg c", code: "c"},
		},
		&GenericError[EDXrr]{
			code: "b",
			err: errors.Join(
				&GenericError[EDXrr]{msg: "msg d", code: "d"},
				&GenericError[EDXrr]{msg: "msg e", code: "e"},
			),
		},
	)
}

// TstTreeCase5 returns a test error tree - case 5.
//
// Shape:
//
//	a─────d────(f)
//	│     │     │
//	│     │   ┌─┴─┐
//	│     │   │   │
//	b     e   g   h
func TstTreeCase5() error {
	fields := map[string]error{
		"f": errors.Join(
			&GenericError[EDXrr]{msg: "msg g", code: "g"},
			&GenericError[EDXrr]{msg: "msg h", code: "h"},
		),
		"a": &GenericError[EDXrr]{
			code: "a",
			err:  &GenericError[EDXrr]{msg: "msg b", code: "b"},
		},
		"d": &GenericError[EDXrr]{
			code: "d",
			err:  &GenericError[EDXrr]{msg: "msg e", code: "e"},
		},
	}
	return NewFields[EDXrr](fields)
}

// TstTreeMeta returns a test error tree with metadata keys. Where the "D"
// metadata key is duplicated in the tree.
//
// Shape:
//
//	     A7,Bb
//	       │
//	  ┌──A6,Cc───┐
//	  │          │
//	A4,Dd    ┌─A5,Ee─┐
//	  │      │       │
//	A3,Ff  A2,Gg   A1,Dh
func TstTreeMeta() error {
	return &GenericError[EDXrr]{
		err: &GenericError[EDXrr]{
			err: errors.Join(
				&GenericError[EDXrr]{
					err:  &GenericError[EDXrr]{msg: "ma3", code: "a3", meta: map[string]any{"A": 3, "F": "f"}},
					meta: map[string]any{"A": 4, "D": "d"},
				},
				&GenericError[EDXrr]{
					err: errors.Join(
						&GenericError[EDXrr]{msg: "ma2", code: "a2", meta: map[string]any{"A": 2, "G": "g"}},
						&GenericError[EDXrr]{msg: "ma1", code: "a1", meta: map[string]any{"A": 1, "D": "h"}},
					),
					meta: map[string]any{"A": 5, "E": "e"},
				},
			),
			meta: map[string]any{"A": 6, "C": "c"},
		},
		meta: map[string]any{"A": 7, "B": "b"},
	}
}
