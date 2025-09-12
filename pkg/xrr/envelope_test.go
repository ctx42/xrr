// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Enclose(t *testing.T) {
	t.Run("nil cause", func(t *testing.T) {
		// --- Given ---
		lead := New("lead", "ECL")

		// --- When ---
		err := Enclose(nil, lead)

		// --- Then ---
		assert.Nil(t, err)
	})

	t.Run("nil lead", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")

		// --- When ---
		err := Enclose(cause, nil)

		// --- Then ---
		assert.NotNil(t, err)
	})

	t.Run("cause and lead provided", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead := New("lead", "ECL")

		// --- When ---
		err := Enclose(cause, lead)

		// --- Then ---
		var enc Envelope
		assert.ErrorAs(t, &enc, err)
		assert.Same(t, cause, enc.cause)
		assert.Same(t, lead, enc.leading)
	})

	t.Run("multiple lead provided", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC")
		lead0 := New("lead0", "ECL")
		lead1 := New("lead1", "ECL")

		// --- When ---
		err := Enclose(cause, lead0, lead1)

		// --- Then ---
		var enc Envelope
		assert.ErrorAs(t, &enc, err)
		assert.Same(t, cause, enc.cause)
		assert.Same(t, lead0, enc.leading)
	})

	t.Run("use passed Envelope - do not nest", func(t *testing.T) {
		// --- Given ---
		evp := Envelope{cause: New("cause", "ECC"), leading: New("lead", "ECL")}

		// --- When ---
		err := Enclose(evp)

		// --- Then ---
		assert.Same(t, evp.cause, err.(Envelope).cause)     // nolint: errorlint
		assert.Same(t, evp.leading, err.(Envelope).leading) // nolint: errorlint
	})

	t.Run("use passed Envelope - override leading", func(t *testing.T) {
		// --- Given ---
		other := New("other", "ECO")
		evp := Envelope{cause: New("cause", "ECC"), leading: New("lead", "ECL")}

		// --- When ---
		err := Enclose(evp, other)

		// --- Then ---
		assert.Same(t, evp.cause, err.(Envelope).cause) // nolint: errorlint
		assert.Same(t, other, err.(Envelope).leading)   // nolint: errorlint
	})
}

func Test_Envelope_Error(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC")
	lead := New("lead", "ECL")
	e := Envelope{cause: cause, leading: lead}

	// --- When ---
	err := e.Error()

	// --- Then ---
	assert.Equal(t, "cause", err)
}

func Test_Envelope_ErrCode(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC")
	lead := New("lead", "ECL")
	e := Envelope{cause: cause, leading: lead}

	// --- When ---
	have := e.ErrCode()

	// --- Then ---
	assert.Equal(t, "ECC", have)
}

func Test_Envelope_Unwrap(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC")
	lead := New("lead", "ECL")
	e := Envelope{cause: cause, leading: lead}

	// --- When ---
	err := e.Unwrap()

	// --- Then ---
	assert.Same(t, cause, err)
}

func Test_Envelope_Lead(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC")
	lead := New("lead", "ECL")
	e := Envelope{cause: cause, leading: lead}

	// --- When ---
	err := e.Lead()

	// --- Then ---
	assert.Same(t, lead, err)
}

func Test_Envelope_Is(t *testing.T) {
	// --- Given ---
	cause := New("cause", "ECC")
	lead := New("lead", "ECL")

	// --- When ---
	err := Envelope{cause: cause, leading: lead}

	// --- Then ---
	assert.ErrorIs(t, lead, err)
	assert.ErrorIs(t, cause, err)
}

func Test_Envelope_MarshalJSON(t *testing.T) {
	t.Run("error with metadata without lead error", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC", Meta().Int("A", 0).Option())

		// --- When ---
		data, err := json.Marshal(Enclose(cause))

		// --- Then ---
		assert.NoError(t, err)
		want := `{"error":"cause","code":"ECC","meta":{"A": 0}}`
		assert.JSON(t, want, string(data))
	})

	t.Run("error with metadata with lead error", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC", Meta().Int("A", 0).Option())
		lead := New("lead", "ECL")

		// --- When ---
		data, err := json.Marshal(Enclose(cause, lead))

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"errors":[
				{"error":"cause","code":"ECC","meta":{"A": 0}}
			]
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("fields error with lead error", func(t *testing.T) {
		// --- Given ---
		cause := Fields{
			"f0": errors.New("f0"),
			"f1": New("f1", "ECF1", Meta().Int("A", 0).Option()),
		}
		lead := New("lead", "ECL")

		// --- When ---
		data, err := json.Marshal(Enclose(cause, lead))

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"fields":{
				"f0":{"error":"f0","code":"ECGeneric"},
				"f1":{"error":"f1","code":"ECF1","meta":{"A": 0}}
			}
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("fields error without lead error", func(t *testing.T) {
		// --- Given ---
		cause := Fields{
			"f0": errors.New("f0"),
			"f1": New("f1", "ECF1", Meta().Int("A", 0).Option()),
		}

		// --- When ---
		data, err := json.Marshal(Enclose(cause))

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"fields error",
			"code":"ECFields",
			"fields":{
				"f0":{"error":"f0","code":"ECGeneric"},
				"f1":{"error":"f1","code":"ECF1","meta":{"A": 0}}
			}
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("single joined error with lead error", func(t *testing.T) {
		// --- Given ---
		cause := errors.Join(errors.New("e0"))
		lead := New("lead", "ECL")

		// --- When ---
		data, err := json.Marshal(Enclose(cause, lead))

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"errors":[
				{"error":"e0","code":"ECGeneric"}
			]
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("joined errors with a lead error", func(t *testing.T) {
		// --- Given ---
		cause := errors.Join(
			New("msg a", "a", Meta().Int("A", 0).Option()),
			errors.New("msg x"),
		)
		lead := New("lead", "ECL")

		// --- When ---
		data, err := json.Marshal(Enclose(cause, lead))

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error": "lead",
			"code": "ECL",
			"errors": [
				{"error":"msg a", "code":"a", "meta": {"A": 0}},
				{"error":"msg x", "code":"ECGeneric"}
			]
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("joined errors without lead error", func(t *testing.T) {
		// --- Given ---
		cause := errors.Join(New("e0", "ECE0", Meta().Int("A", 0).Option()), errors.New("e1"))

		// --- When ---
		data, err := json.Marshal(Enclose(cause))

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"e0",
			"code":"ECE0",
			"errors":[
				{"error": "e1", "code": "ECGeneric"}
			],
			"meta": {"A": 0}
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("marshal error", func(t *testing.T) {
		// --- Given ---
		e1 := &TErrMarshalJSON{New("msg a", "a")}
		cause := errors.Join(New("msg b", "b"), e1)

		// --- When ---
		data, err := json.Marshal(Enclose(cause))

		// --- Then ---
		var jme *json.MarshalerError
		assert.ErrorAs(t, &jme, err)
		want := "json: error calling MarshalJSON for type xrr.Envelope: " +
			"json: error calling MarshalJSON for type *xrr.TErrMarshalJSON: msg a"
		assert.Equal(t, want, jme.Error())
		assert.Nil(t, data)
	})
}

func Test_encloseFieldsError(t *testing.T) {
	t.Run("lead without metadata", func(t *testing.T) {
		// --- Given ---
		cause := Fields{
			"f0": errors.New("f0"),
			"f1": New("f1", "ECF1", Meta().Int("A", 0).Option()),
		}
		lead := New("lead", "ECL")

		// --- When ---
		data, err := encloseFieldsError(lead, cause)

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"fields":{
				"f0":{"error":"f0","code":"ECGeneric"},
				"f1":{"error":"f1","code":"ECF1","meta":{"A": 0}}
			}
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("lead with metadata", func(t *testing.T) {
		// --- Given ---
		cause := Fields{
			"f0": errors.New("f0"),
			"f1": New("f1", "ECF1", Meta().Int("A", 0).Option()),
		}
		lead := New("lead", "ECL", Meta().Int("B", 1).Option())

		// --- When ---
		data, err := encloseFieldsError(lead, cause)

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"fields":{
				"f0":{"error":"f0","code":"ECGeneric"},
				"f1":{"error":"f1","code":"ECF1","meta":{"A": 0}}
			},
			"meta":{"B": 1}
		}`
		assert.JSON(t, want, string(data))
	})
}

func Test_encloseMultiError(t *testing.T) {
	t.Run("lead without metadata", func(t *testing.T) {
		// --- Given ---
		e0 := New("e0", "ECE0", Meta().Int("A", 0).Option())
		e1 := errors.New("e1")
		lead := New("lead", "ECL")

		// --- When ---
		data, err := encloseMultiError(lead, e0, e1)

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"errors":[
				{"error":"e0","code":"ECE0","meta":{"A": 0}},
				{"error":"e1","code":"ECGeneric"}
			]
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("lead with metadata", func(t *testing.T) {
		// --- Given ---
		cause := New("cause", "ECC", Meta().Int("A", 0).Option())
		lead := New("lead", "ECL", Meta().Int("B", 1).Option())

		// --- When ---
		data, err := encloseMultiError(lead, cause)

		// --- Then ---
		assert.NoError(t, err)
		want := `
		{
			"error":"lead",
			"code":"ECL",
			"errors":[
				{"error":"cause","code":"ECC","meta":{"A": 0}}
			],
			"meta":{"B": 1} 
		}`
		assert.JSON(t, want, string(data))
	})

	t.Run("lead with no errors", func(t *testing.T) {
		// --- Given ---
		lead := New("lead", "ECL", Meta().Int("A", 0).Option())

		// --- When ---
		data, err := encloseMultiError(lead)

		// --- Then ---
		assert.NoError(t, err)
		want := `{"error":"lead", "code":"ECL", "meta":{"A": 0}}`
		assert.JSON(t, want, string(data))
	})
}

func Test_marshalError(t *testing.T) {
	t.Run("std error", func(t *testing.T) {
		// --- Given ---
		e := errors.New("e")

		// --- When ---
		data, err := marshalError(e)

		// --- Then ---
		assert.NoError(t, err)
		assert.JSON(t, `{"error": "e", "code": "ECGeneric"}`, string(data))
	})

	t.Run("xrr error", func(t *testing.T) {
		// --- Given ---
		e := New("msg a", "a")

		// --- When ---
		data, err := marshalError(e)

		// --- Then ---
		assert.NoError(t, err)
		assert.JSON(t, `{"error": "msg a", "code": "a"}`, string(data))
	})

	t.Run("marshal error", func(t *testing.T) {
		// --- Given ---
		e := &TErrMarshalJSON{errors.New("e")}

		// --- When ---
		data, err := marshalError(e)

		// --- Then ---
		wMsg := "json: " +
			"error calling MarshalJSON for type *xrr.TErrMarshalJSON: e"
		assert.ErrorEqual(t, wMsg, err)
		assert.Nil(t, data)
	})
}
