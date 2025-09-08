package xrr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_GetCode_tabular(t *testing.T) {
	var err0 error
	var err1 *Error

	tt := []struct {
		testN string

		exp string
		err error
	}{
		{"with code", "ECode", New("msg", "ECode")},
		{"simple error", ECGeneric, errors.New("msg")},
		{"nil", "", nil},
		{"nil error", "", err0},
		{"nil Error", "", err1},

		{
			"the first error code from joined errors is returned case 1",
			"ECode",
			fmt.Errorf("%w: %w", New("msg0", "ECode"), errors.New("msg1")),
		},
		{
			"the first error code from joined errors is returned case 2",
			"ECode0",
			fmt.Errorf("%w: %w", New("msg0", "ECode0"), New("msg1", "ECode1")),
		},
		{
			"the first error code from joined errors is returned case 3",
			ECGeneric,
			fmt.Errorf("%w: %w", errors.New("msg0"), New("msg1", "ECode")),
		},
		{
			"the first error code from joined errors is returned case 1",
			"ECode",
			errors.Join(New("msg0", "ECode"), errors.New("msg1")),
		},
		{
			"the first error code from joined errors is returned case 2",
			"ECode0",
			errors.Join(New("msg0", "ECode0"), New("msg1", "ECode1")),
		},
		{
			"the first error code from joined errors is returned case 3",
			ECGeneric,
			errors.Join(errors.New("msg0"), New("msg1", "ECode")),
		},
		{
			"error wrapped once",
			"ECode",
			fmt.Errorf("comment: %w", New("msg1", "ECode")),
		},
		{
			"error wrapped many times",
			"ECode",
			fmt.Errorf(
				"second: %w",
				fmt.Errorf("first: %w", New("msg1", "ECode")),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := GetCode(tc.err)

			// --- Then ---
			assert.Equal(t, tc.exp, have)
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
		{"std error", errors.New("msg"), nil},
		{
			"xrr std wrapped with metadata",
			fmt.Errorf("wrapped: %w", New("msg", "ECode")),
			[]string{"ECode"},
		},
		{
			"xrr std wrapped multiple times",
			fmt.Errorf("2: %w", fmt.Errorf("1: %w", New("msg", "ECode"))),
			[]string{"ECode"},
		},
		{
			"joined errors",
			errors.Join(New("msg0", "ECode0"), New("msg1", "ECode1")),
			[]string{"ECode0", "ECode1"},
		},
		{
			"joined and std wrapped errors",
			errors.Join(
				New("msg0", "ECode0"),
				New("msg1", "ECode1"),
				fmt.Errorf("wrapped: %w", New("msg2", "ECode2")),
			),
			[]string{"ECode0", "ECode1", "ECode2"},
		},
		{
			"tree",
			&Error{
				error: &Error{
					error: errors.Join(
						&Error{error: New("msg3", "EC3"), code: "EC4"},
						&Error{
							error: errors.Join(
								New("msg1", "EC1"),
								New("msg2", "EC2"),
							),
							code: "EC5",
						},
					),
					code: "EC6",
				},
				code: "EC7",
			},
			[]string{"EC1", "EC2", "EC3", "EC4", "EC5", "EC6", "EC7"},
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
		{"std error", errors.New("msg"), nil},
		{"xrr without metadata", New("msg", "ECode"), nil},
		{
			"xrr with metadata",
			New("msg", "ECode", Meta().Int("A", 1).Bool("B", true).Option()),
			map[string]any{"A": 1, "B": true},
		},
		{
			"xrr std wrapped with metadata",
			fmt.Errorf(
				"wrapped: %w",
				New("msg", "ECode", Meta().Int("A", 1).Bool("B", true).Option()),
			),
			map[string]any{"A": 1, "B": true},
		},
		{
			"xrr std wrapped multiple times with metadata",
			fmt.Errorf(
				"second: %w",
				fmt.Errorf(
					"first: %w",
					New("msg", "ECode", Meta().Int("A", 1).Bool("B", true).Option()),
				),
			),
			map[string]any{"A": 1, "B": true},
		},
		{
			"joined errors",
			errors.Join(
				New("msg", "ECode", Meta().Int("A", 1).Int("B", 1).Option()),
				New("msg", "ECode", Meta().Int("A", 1).Int("B", 2).Option()),
			),
			map[string]any{"A": 1, "B": 1},
		},
		{
			"joined and std wrapped errors",
			errors.Join(
				New("msg", "ECode", Meta().Int("A", 1).Int("B", 1).Option()),
				New("msg", "ECode", Meta().Int("A", 2).Int("B", 2).Option()),
				fmt.Errorf(
					"wrapped: %w",
					New("msg", "ECode", Meta().Int("A", 3).Int("C", 1).Option()),
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
				error: &Error{
					error: errors.Join(
						&Error{
							error: &Error{meta: map[string]any{"A": 3, "F": 3}},
							meta:  map[string]any{"A": 4, "D": 2},
						},
						&Error{
							error: errors.Join(
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
