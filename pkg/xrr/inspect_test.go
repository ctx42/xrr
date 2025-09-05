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
