package xrr

import (
	"errors"
	"fmt"
	"testing"

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
	}

	// TODO(rz): nil errors are ignored - implement Fields first.

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
		//        │
		//   ┌──A6,C1───┐
		//   │          │
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
