// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_IsCode_tabular(t *testing.T) {
	var err0 error
	var err1 *Error

	tt := []struct {
		testN string

		err  error
		code string
		want bool
	}{
		{"simple has code", New("msg a", "a"), "a", true},
		{"simple has different code", New("msg b", "b"), "a", false},
		{"nil", nil, "a", false},
		{"nil error", err0, "a", false},
		{"nil Error", err1, "a", false},
		{"std error", errors.New("msg x"), "a", false},
		{
			"joined errors with searched code",
			fmt.Errorf("%w: %w", errors.New("msg x"), New("msg a", "a")),
			"a",
			true,
		},
		{
			"joined errors without searched code",
			fmt.Errorf("%w: %w", New("msg a", "a"), New("msg b", "b")),
			"x",
			false,
		},
		{
			"wrapped error with searched code",
			fmt.Errorf("comment: %w", New("msg a", "a")),
			"a",
			true,
		},
		{
			"multiple wrapped error with searched code",
			fmt.Errorf("1: %w", fmt.Errorf("2: %w", New("msg a", "a"))),
			"a",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := IsCode(tc.err, tc.code)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_GetCode_tabular(t *testing.T) {
	var err0 error
	var err1 *Error

	tt := []struct {
		testN string

		err  error
		want string
	}{
		{"with code", New("msg a", "a"), "a"},
		{"simple error", errors.New("msg x"), ECGeneric},
		{"nil", nil, ""},
		{"nil error", err0, ""},
		{"nil Error", err1, ""},
		{
			"wrapped error",
			fmt.Errorf("1: %w", New("msg a", "a")),
			ECGeneric,
		},
		{
			"wrapped errs",
			fmt.Errorf("%w: %w", New("msg a", "a"), errors.New("msg x")),
			ECGeneric,
		},
		{
			"error wrapped with xrr.Wrap",
			Wrap(New("msg a", "a"), Meta().Int("A", 1).Option()),
			"a",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := GetCode(tc.err)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_GetCodes_tabular(t *testing.T) {
	var err error

	tt := []struct {
		testN string

		err  error
		want []string
	}{
		{"nil", nil, nil},
		{"nil typed", err, nil},
		{"std error", errors.New("msg x"), []string{ECGeneric}},
		{
			"xrr std wrapped with metadata",
			fmt.Errorf("wrapped: %w", New("msg a", "a")),
			[]string{ECGeneric, "a"},
		},
		{
			"xrr std wrapped multiple times",
			fmt.Errorf("2: %w", fmt.Errorf("1: %w", New("msg a", "a"))),
			[]string{ECGeneric, "a"},
		},
		{
			"joined errors",
			errors.Join(New("msg a", "a"), New("msg b", "b")),
			[]string{"a", "b"},
		},
		{
			"joined and std wrapped errors",
			errors.Join(
				New("msg a", "a"),
				New("msg b", "b"),
				fmt.Errorf("wrapped: %w", New("msg c", "c")),
			),
			[]string{"a", "b", ECGeneric, "c"},
		},
		{
			"tree",
			TstTreeCase1(),
			[]string{"a", "b", "c", "e", "d", "f", "g"},
		},
		{
			"does not return repeated codes",
			&Error{
				code: "a",
				err: &Error{
					code: "b",
					err: &Error{
						code: "a",
						err: &Error{
							code: "c",
						},
					},
				},
			},
			[]string{"a", "b", "c"},
		},
		{
			"errors without error code",
			&Error{
				code: "a",
				err: &Error{
					code: "b",
					err:  errors.New("msg x"),
				},
			},
			[]string{"a", "b", ECGeneric},
		},
		{
			"nil errors are ignored",
			&Error{
				code: "a",
				err: &Error{
					code: "b",
					err: &Fields{
						"x": nil,
						"y": New("msg y", "y"),
						"z": nil,
					},
				},
			},
			[]string{"a", "b", "y"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := GetCodes(tc.err)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_GetMeta_tabular(t *testing.T) {
	var err error

	tt := []struct {
		testN string

		err  error
		want map[string]any
	}{
		{"nil", nil, nil},
		{"nil typed", err, nil},
		{"std error", errors.New("msg x"), nil},
		{"xrr without metadata", New("msg a", "a"), nil},
		{
			"xrr with metadata",
			New("msg a", "a", Meta().Int("A", 1).Bool("B", true).Option()),
			map[string]any{"A": 1, "B": true},
		},
		{
			"xrr std wrapped with metadata",
			fmt.Errorf(
				"wrapped: %w",
				New("msg a", "a", Meta().Int("A", 1).Bool("B", true).Option()),
			),
			map[string]any{"A": 1, "B": true},
		},
		{
			"xrr std wrapped multiple times with metadata",
			fmt.Errorf(
				"second: %w",
				fmt.Errorf(
					"first: %w",
					New("msg a", "a", Meta().Int("A", 1).Bool("B", true).Option()),
				),
			),
			map[string]any{"A": 1, "B": true},
		},
		{
			"joined errors",
			errors.Join(
				New("msg a", "a", Meta().Int("A", 1).Int("B", 1).Option()),
				New("msg b", "b", Meta().Int("A", 1).Int("B", 2).Option()),
			),
			map[string]any{"A": 1, "B": 1},
		},
		{
			"joined and std wrapped errors",
			errors.Join(
				New("msg a", "a", Meta().Int("A", 1).Int("B", 1).Option()),
				New("msg b", "b", Meta().Int("A", 2).Int("B", 2).Option()),
				fmt.Errorf(
					"wrapped: %w",
					New("msg c", "c", Meta().Int("A", 3).Int("C", 1).Option()),
				),
			),
			map[string]any{"A": 1, "B": 1, "C": 1},
		},
		//      A7,B0
		//        │
		//   ┌──A6,C1───┐
		//   │          │
		// A4,D2    ┌─A5,E2─┐
		//   │      │       │
		// A3,F3  A2,G3   A1,H3
		//
		{
			"tree",
			&Error{
				err: &Error{
					err: errors.Join(
						&Error{
							err:  &Error{meta: map[string]any{"A": 3, "F": 3}},
							meta: map[string]any{"A": 4, "D": 2},
						},
						&Error{
							err: errors.Join(
								&Error{meta: map[string]any{"A": 2, "G": 3}},
								&Error{meta: map[string]any{"A": 1, "H": 3}},
							),
							meta: map[string]any{"A": 5, "E": 2},
						},
					),
					meta: map[string]any{"A": 6, "C": 1},
				},
				meta: map[string]any{"A": 7, "B": 0},
			},
			map[string]any{
				"A": 7,
				"B": 0,
				"C": 1,
				"D": 2,
				"E": 2,
				"F": 3,
				"G": 3,
				"H": 3,
			},
		},
		{
			"fields",
			&Fields{
				"a": New("msg a", "a", Meta().Int("A", 1).Int("B", 1).Option()),
				"b": New("msg b", "b", Meta().Int("A", 1).Int("B", 2).Option()),
			},
			map[string]any{"A": 1, "B": 1},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := GetMeta(tc.err)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_GetBool(t *testing.T) {
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": true},
			err: &Error{
				meta: map[string]any{"A": false, "B": 3},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value bool
		exist bool
	}{
		{"nil error", nil, "key", false, false},
		{"not existing key", tree(), "X", false, false},
		{"returns the first found key", tree(), "A", true, true},
		{"key found but type mismatch", tree(), "B", false, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetBool(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetStr(t *testing.T) {
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": "1"},
			err: &Error{
				meta: map[string]any{"A": "2", "B": 3},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value string
		exist bool
	}{
		{"nil error", nil, "key", "", false},
		{"not existing key", tree(), "X", "", false},
		{"returns the first found key", tree(), "A", "1", true},
		{"key found but type mismatch", tree(), "B", "", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetStr(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetInt(t *testing.T) {
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": 1},
			err: &Error{
				meta: map[string]any{"A": 2, "B": "3"},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value int
		exist bool
	}{
		{"nil error", nil, "key", 0, false},
		{"not existing key", tree(), "X", 0, false},
		{"returns the first found key", tree(), "A", 1, true},
		{"key found but type mismatch", tree(), "B", 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetInt(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetInt64(t *testing.T) {
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": int64(1)},
			err: &Error{
				meta: map[string]any{"A": int64(2), "B": "3"},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value int64
		exist bool
	}{
		{"nil error", nil, "key", 0, false},
		{"not existing key", tree(), "X", 0, false},
		{"returns the first found key", tree(), "A", 1, true},
		{"key found but type mismatch", tree(), "B", 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetInt64(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetFloat64(t *testing.T) {
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": float64(1)},
			err: &Error{
				meta: map[string]any{"A": float64(2), "B": "3"},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value float64
		exist bool
	}{
		{"nil error", nil, "key", 0, false},
		{"not existing key", tree(), "X", 0, false},
		{"returns the first found key", tree(), "A", 1, true},
		{"key found but type mismatch", tree(), "B", 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetFloat64(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetTime(t *testing.T) {
	tim1 := time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC)
	tim2 := time.Date(2002, 2, 2, 2, 2, 2, 0, time.UTC)
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": tim1},
			err: &Error{
				meta: map[string]any{"A": tim2, "B": "3"},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value time.Time
		exist bool
	}{
		{"nil error", nil, "key", time.Time{}, false},
		{"not existing key", tree(), "X", time.Time{}, false},
		{"returns the first found key", tree(), "A", tim1, true},
		{"key found but type mismatch", tree(), "B", time.Time{}, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetTime(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetDuration(t *testing.T) {
	tree := func() error {
		return &Error{
			meta: map[string]any{"A": time.Second},
			err: &Error{
				meta: map[string]any{"A": time.Hour, "B": "3"},
			},
		}
	}

	tt := []struct {
		testN string

		err   error
		key   string
		value time.Duration
		exist bool
	}{
		{"nil error", nil, "key", 0, false},
		{"not existing key", tree(), "X", 0, false},
		{"returns the first found key", tree(), "A", time.Second, true},
		{"key found but type mismatch", tree(), "B", 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetDuration(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_getKey(t *testing.T) {
	tt := []struct {
		testN string

		err   error
		key   string
		value string
		exist bool
	}{
		{"nil error", nil, "A", "", false},
		{"std error", errors.New("a"), "A", "", false},
		{"no key simple", New("a", "eca"), "A", "", false},
		{"key fund but type mismatch", TstTreeMeta(), "A", "", false},
		{"returns the first found key", TstTreeMeta(), "D", "d", true},
		{"leaf node key", TstTreeMeta(), "G", "g", true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := getKey[string](tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_walk(t *testing.T) {
	t.Run("tree configuration 1", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase1()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "abcedfg", have)
	})

	t.Run("tree configuration 2", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase2()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "abcehdfgi", have)
	})

	t.Run("tree configuration 3", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "abcehdfig", have)
	})

	t.Run("tree configuration 4", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase4()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "acbde", have)
	})

	t.Run("tree configuration 5", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase5()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "abdegh", have)
	})

	t.Run("stop after the first one", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		var have string
		cb := func(err error) bool { have += GetCode(err); return false }

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "a", have)
	})

	t.Run("stop walking after first", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		var have string
		cb := func(err error) bool {
			have += GetCode(err)
			return len(have) < 3
		}

		// --- When ---
		walk(e, cb)

		// --- Then ---
		assert.Equal(t, "abc", have)
	})
}

func Test_walkRev(t *testing.T) {
	t.Run("tree configuration 1", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase1()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "gfdecba", have)
	})

	t.Run("tree configuration 2", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase2()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "igfdhecba", have)
	})

	t.Run("tree configuration 3", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "gifdhecba", have)
	})

	t.Run("tree configuration 4", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase4()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "edbca", have)
	})

	t.Run("tree configuration 5", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase5()

		var have string
		cb := func(err error) bool { have += GetCode(err); return true }

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "hgedba", have)
	})

	t.Run("stop after the first one", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		var have string
		cb := func(err error) bool { have += GetCode(err); return false }

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "g", have)
	})

	t.Run("stop walking after first", func(t *testing.T) {
		// --- Given ---
		e := TstTreeCase3()

		var have string
		cb := func(err error) bool {
			have += GetCode(err)
			return len(have) < 3
		}

		// --- When ---
		walkReverse(e, cb)

		// --- Then ---
		assert.Equal(t, "gif", have)
	})
}
