// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_Mask(t *testing.T) {
	t.Run("nil cause returns nil", func(t *testing.T) {
		// --- When ---
		err := Mask(nil, New("lead", "ECL"))

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("nil lead returns nil", func(t *testing.T) {
		// --- When ---
		err := Mask(New("cause", "ECC"), nil)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("cause and lead provided", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := Mask(cause, lead)

		// --- Then ---
		var enc Masked
		assert.ErrorAs(t, &enc, err)
		assert.Same(t, cause, enc.cause)
		assert.Same(t, lead, enc.lead)
	})
}

func Test_Masked_Error(t *testing.T) {
	// --- When ---
	err := Mask(New("cause", "ECC"), New("lead", "ECL"))

	// --- Then ---
	assert.Equal(t, "lead", err.Error())
}

func Test_Masked_ErrorCode(t *testing.T) {
	// --- When ---
	err := Mask(New("cause", "ECC"), New("lead", "ECL"))

	// --- Then ---
	assert.Equal(t, "ECL", GetCode(err))
}

func Test_Masked_Cause(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC")

	// --- When ---
	err := Mask(cause, New("lead", "ECL"))

	// --- Then ---
	var me Masked
	assert.ErrorAs(t, &me, err)
	assert.Same(t, cause, me.Cause())
}

func Test_Masked_Lead(t *testing.T) {
	// --- Given ---
	lead := New("lead", "ECL")

	// --- When ---
	err := Mask(New("cause", "ECC"), lead)

	// --- Then ---
	var me Masked
	assert.ErrorAs(t, &me, err)
	assert.Same(t, lead, me.Lead())
}

func Test_Masked_MarshalJSON(t *testing.T) {
	t.Run("lead without metadata", func(t *testing.T) {
		// --- Given ---
		e := Mask(New("cause", "ECC"), New("lead", "ECL"))

		// --- When ---
		data := must.Value(json.Marshal(e))

		// --- Then ---
		want := `{"error":"lead","code":"ECL"}`
		assert.JSON(t, want, string(data))
	})

	t.Run("lead with metadata", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC", Meta().Int("X", 1).Option())
		lead := New("lead", "ECL", Meta().Int("Y", 2).Option())
		e := Mask(cause, lead)

		// --- When ---
		data := must.Value(json.Marshal(e))

		// --- Then ---
		want := `{"error":"lead","code":"ECL","meta":{"Y":2}}`
		assert.JSON(t, want, string(data))
	})
}

func Test_Masked_Unwrap(t *testing.T) {
	t.Run("returns cause and lead", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := Mask(cause, lead)

		// --- Then ---
		var me Masked
		assert.ErrorAs(t, &me, err)
		errs := me.Unwrap()
		assert.Equal(t, 2, len(errs))
		assert.Same(t, cause, errs[0])
		assert.Same(t, lead, errs[1])
	})

	t.Run("Is finds both cause and lead in chain", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := Mask(cause, lead)

		// --- Then ---
		assert.ErrorIs(t, cause, err)
		assert.ErrorIs(t, lead, err)
	})
}
