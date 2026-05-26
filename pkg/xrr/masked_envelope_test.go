// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_EncloseMasked(t *testing.T) {
	t.Run("nil cause returns nil", func(t *testing.T) {
		// --- When ---
		err := EncloseMasked(nil, New("lead", "ECL"))

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("nil lead returns nil", func(t *testing.T) {
		// --- When ---
		err := EncloseMasked(New("cause", "ECC"), nil)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("cause and lead provided", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := EncloseMasked(cause, lead)

		// --- Then ---
		var enc MaskedEnvelope
		assert.ErrorAs(t, &enc, err)
		assert.Same(t, cause, enc.cause)
		assert.Same(t, lead, enc.lead)
	})

	t.Run("Is finds both cause and lead in chain", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := EncloseMasked(cause, lead)

		// --- Then ---
		assert.ErrorIs(t, cause, err)
		assert.ErrorIs(t, lead, err)
	})

	t.Run("Unwrap returns cause and lead", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := EncloseMasked(cause, lead)

		// --- Then ---
		var me MaskedEnvelope
		assert.ErrorAs(t, &me, err)
		errs := me.Unwrap()
		assert.Equal(t, 2, len(errs))
		assert.Same(t, cause, errs[0])
		assert.Same(t, lead, errs[1])
	})
}

func Test_MaskedEnvelope_Error(t *testing.T) {
	// --- Given ---
	err := EncloseMasked(New("cause", "ECC"), New("lead", "ECL"))

	// --- Then ---
	assert.Equal(t, "lead", err.Error())
}

func Test_MaskedEnvelope_ErrorCode(t *testing.T) {
	// --- Given ---
	err := EncloseMasked(New("cause", "ECC"), New("lead", "ECL"))

	// --- Then ---
	assert.Equal(t, "ECL", GetCode(err))
}

func Test_MaskedEnvelope_MarshalJSON(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC", Meta().Int("X", 1).Option())
	lead := New("lead", "ECL", Meta().Int("Y", 2).Option())
	e := EncloseMasked(cause, lead)

	// --- Then ---
	data := must.Value(json.Marshal(e))
	want := `{"error":"lead","code":"ECL","meta":{"Y":2}}`
	assert.JSON(t, want, string(data))
}
