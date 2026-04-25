// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_NewDomainFields(t *testing.T) {
	t.Run("stores map directly", func(t *testing.T) {
		// --- Given ---
		m := map[string]error{"f0": ErrTst}

		// --- When ---
		have := NewDomainFields[EDGeneric](m)

		// --- Then ---
		assert.Same(t, ErrTst, have.fields["f0"])
		assert.Len(t, 1, have.fields)
		// Map is aliased, not copied.
		m["f1"] = errors.New("extra")
		assert.Len(t, 2, have.fields)
	})

	t.Run("nil map", func(t *testing.T) {
		// --- When ---
		have := NewDomainFields[EDGeneric](nil)

		// --- Then ---
		assert.Nil(t, have.fields)
	})
}

func Test_FieldsFactory(t *testing.T) {
	t.Run("create error", func(t *testing.T) {
		// --- Given ---
		have := FieldsFactory[EDGeneric]()

		// --- When ---
		err := have("field", ErrTst)

		// --- Then ---
		fe, _ := assert.SameType(t, &GenericFields[EDGeneric]{}, err)
		e, _ := assert.HasKey(t, "field", fe.ErrorFields())
		assert.Same(t, ErrTst, e)
	})

	t.Run("nil error returns nil", func(t *testing.T) {
		// --- Given ---
		have := FieldsFactory[EDGeneric]()

		// --- When ---
		err := have("field", nil)

		// --- Then ---
		assert.Nil(t, err)
	})
}

func Test_GetFields(t *testing.T) {
	// --- Given ---
	err := New("message", "code")
	fls := TErrorFields(map[string]error{"key": err})

	// --- When ---
	m := GetFields(fls)

	// --- Then ---
	assert.NotNil(t, m)
	assert.Len(t, 1, m)
	_, _ = assert.HasKey(t, "key", m)
	assert.Same(t, err, m["key"])
}

func Test_GetFieldError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		have := GetFieldError(nil, "f0")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("not instance of Fields", func(t *testing.T) {
		// --- Given ---
		err := errors.New("em0")

		// --- When ---
		have := GetFieldError(err, "f0")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("field does not exist", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
			},
		}

		// --- When ---
		have := GetFieldError(fs, "f1")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("field exists", func(t *testing.T) {
		// --- Given ---
		err := errors.New("em0")
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": err,
			},
		}

		// --- When ---
		have := GetFieldError(fs, "f0")

		// --- Then ---
		assert.Same(t, err, have)
	})
}

func Test_FieldErrorIs(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		// --- When ---
		have := FieldErrorIs(nil, "f0", io.EOF)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not instance of Fields", func(t *testing.T) {
		// --- Given ---
		err := errors.New("em0")

		// --- When ---
		have := FieldErrorIs(err, "f0", io.EOF)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("field does not exist", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
			},
		}

		// --- When ---
		have := FieldErrorIs(fs, "f1", io.EOF)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("field exist", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": io.EOF,
			},
		}

		// --- When ---
		have := FieldErrorIs(fs, "f0", io.EOF)

		// --- Then ---
		assert.True(t, have)
	})
}

func Test_FieldNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": errors.New("em1"),
			},
		}

		// --- When ---
		have := FieldNames(fs)

		// --- Then ---
		assert.Equal(t, []string{"f0", "f1"}, have)
	})

	t.Run("empty", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{}

		// --- When ---
		have := FieldNames(fs)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("nil", func(t *testing.T) {
		// --- When ---
		have := FieldNames(nil)

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_MergeFields(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// --- Given ---
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
			},
		}
		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f1": errors.New("m1"),
				"f2": errors.New("m2"),
			},
		}

		// --- When ---
		have := MergeFields[EDGeneric](fs0, fs1)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": errors.New("m1"),
				"f2": errors.New("m2"),
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("later keys overwrite previous ones", func(t *testing.T) {
		// --- Given ---
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
			},
		}
		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f1": errors.New("m1"),
				"f0": errors.New("m2"),
			},
		}

		// --- When ---
		have := MergeFields[EDGeneric](fs0, fs1)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m2"),
				"f1": errors.New("m1"),
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("dop not override errors with later nil errors", func(t *testing.T) {
		// --- Given ---
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
			},
		}
		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f1": errors.New("m1"),
				"f0": nil,
			},
		}

		// --- When ---
		have := MergeFields[EDGeneric](fs0, fs1)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": errors.New("m1"),
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("nil errors are not skipped", func(t *testing.T) {
		// --- Given ---
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": nil,
			},
		}
		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f2": errors.New("m2"),
				"f3": errors.New("m3"),
				"f4": nil,
			},
		}

		// --- When ---
		have := MergeFields[EDGeneric](fs0, fs1)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": nil,
				"f2": errors.New("m2"),
				"f3": errors.New("m3"),
				"f4": nil,
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("nil fields are skipped", func(t *testing.T) {
		// --- Given ---
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": nil,
			},
		}
		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f2": errors.New("m2"),
				"f3": errors.New("m3"),
				"f4": nil,
			},
		}

		// --- When ---
		have := MergeFields[EDGeneric](fs0, nil, fs1, nil)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": nil,
				"f2": errors.New("m2"),
				"f3": errors.New("m3"),
				"f4": nil,
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("call with an empty slice returns nil", func(t *testing.T) {
		// --- When ---
		have := MergeFields[EDGeneric]()

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("call with nil", func(t *testing.T) {
		// --- When ---
		have := MergeFields[EDGeneric](nil)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("call with multiple nils", func(t *testing.T) {
		// --- When ---
		have := MergeFields[EDGeneric](nil, nil, nil)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("not Field instances get fake indexed field names", func(t *testing.T) {
		// --- Given ---
		e0 := errors.New("abc")
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("m0"),
				"f1": nil,
			},
		}
		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f2": errors.New("m2"),
				"f3": errors.New("m3"),
				"f4": nil,
			},
		}
		e3 := errors.New("def")

		// --- When ---
		have := MergeFields[EDGeneric](e0, fs0, fs1, e3)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"__field__0": errors.New("abc"),
				"__field__3": errors.New("def"),
				"f0":         errors.New("m0"),
				"f1":         nil,
				"f2":         errors.New("m2"),
				"f3":         errors.New("m3"),
				"f4":         nil,
			},
		}
		assert.Equal(t, want, have)
	})
}

func Test_GenericFields_ErrorFields(t *testing.T) {
	// --- Given ---
	fields := map[string]error{
		"f0": errors.New("em0"),
		"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
		"f2": New("em2", "ECode2"),
	}
	fe := &GenericFields[EDGeneric]{fields: fields}

	// --- When ---
	have := fe.ErrorFields()

	// --- Then ---
	assert.Equal(t, fields, have)
}

func Test_GenericFields_Error(t *testing.T) {
	t.Run("not nested", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fs.Error()

		// --- Then ---
		want := "" +
			"f0: em0; " +
			"f1: em1; " +
			"f2: em2"
		assert.Equal(t, want, have)
	})

	t.Run("nested", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fs.Error()

		// --- Then ---
		want := "" +
			"f0.s0: em00; " +
			"f0.s1: em01; " +
			"f0.s2.s0: em020; " +
			"f1: em1; " +
			"f2: em2"
		assert.Equal(t, want, have)
	})

	t.Run("nil error", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": nil,
			},
		}

		// --- When ---
		have := fs.Error()

		// --- Then ---
		assert.Empty(t, have)
	})

	t.Run("the first sorted key is nil", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"a": nil,
				"b": errors.New("em0"),
			},
		}

		// --- When ---
		have := fs.Error()

		// --- Then ---
		assert.Equal(t, "b: em0", have)
	})
}

func Test_GenericFields_Unwrap(t *testing.T) {
	t.Run("not nested", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": errors.New("em1"),
				"f2": errors.New("em2"),
				"f3": nil,
			},
		}

		// --- When ---
		have := fs.Unwrap()

		// --- Then ---
		assert.Len(t, 3, have)
		assert.ErrorEqual(t, "f0: em0", have[0])
		assert.ErrorEqual(t, "f1: em1", have[1])
		assert.ErrorEqual(t, "f2: em2", have[2])
	})

	t.Run("nested", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": errors.New("em10"),
						"s1": errors.New("em11"),
						"s2": nil,
					},
				},
				"f2": errors.New("em2"),
				"f3": nil,
			},
		}

		// --- When ---
		have := fs.Unwrap()

		// --- Then ---
		assert.Len(t, 4, have)
		assert.ErrorEqual(t, "f0: em0", have[0])
		assert.ErrorEqual(t, "f1.s0: em10", have[1])
		assert.ErrorEqual(t, "f1.s1: em11", have[2])
		assert.ErrorEqual(t, "f2: em2", have[3])
	})
}

func Test_GenericFields_Is(t *testing.T) {
	t.Run("true not nested fields", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": ErrTst,
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fs.Is(ErrTst)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("works on all levels", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": ErrTst,
							},
						},
					},
				},
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fs.Is(ErrTst)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("nil field error", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": nil,
			},
		}

		// --- When ---
		have := fs.Is(ErrTst)

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("nil other error", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
			},
		}

		// --- When ---
		have := fs.Is(nil)

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_GenericFields_Format(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{}

		// --- When ---
		have := fmt.Sprintf("%s", fs)

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("not nested s", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%s", fs)

		// --- Then ---
		want := "f0: em0; f1: em1; f2: em2"
		assert.Equal(t, want, have)
	})

	t.Run("not nested q", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%q", fs)

		// --- Then ---
		want := `"f0: em0; f1: em1; f2: em2"`
		assert.Equal(t, want, have)
	})

	t.Run("not nested v", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%v", fs)

		// --- Then ---
		want := "f0: em0; f1: em1; f2: em2"
		assert.Equal(t, want, have)
	})

	t.Run("not nested plus v", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": errors.New("em0"),
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%+v", fs)

		// --- Then ---
		want := "f0: em0 (ECGeneric); f1: em1 (ECode1); f2: em2 (ECode2)"
		assert.Equal(t, want, have)
	})

	t.Run("nested s", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%s", fs)

		// --- Then ---
		want := "" +
			"f0.s0: em00; " +
			"f0.s1: em01; " +
			"f0.s2.s0: em020; " +
			"f1: em1; " +
			"f2: em2"
		assert.Equal(t, want, have)
	})

	t.Run("nested q", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%q", fs)

		// --- Then ---
		want := `"f0.s0: em00; f0.s1: em01; f0.s2.s0: em020; f1: em1; f2: em2"`
		assert.Equal(t, want, have)
	})

	t.Run("nested v", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%v", fs)

		// --- Then ---
		want := "f0.s0: em00; f0.s1: em01; f0.s2.s0: em020; f1: em1; f2: em2"
		assert.Equal(t, want, have)
	})

	t.Run("nested plus v", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have := fmt.Sprintf("%+v", fs)

		// --- Then ---
		want := "" +
			"f0.s0: em00 (ECode00); " +
			"f0.s1: em01 (ECode01); " +
			"f0.s2.s0: em020 (ECGeneric); " +
			"f1: em1 (ECode1); " +
			"f2: em2 (ECode2)"
		assert.Equal(t, want, have)
	})
}

func Test_GenericFields_Flatten(t *testing.T) {
	t.Run("flatten single", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": errors.New("em00"),
						"s1": errors.New("em01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1"),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		err := fs.Flatten()

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0.s0":    errors.New("em00"),
				"f0.s1":    errors.New("em01"),
				"f0.s2.s0": errors.New("em020"),
				"f1":       New("em1", "ECode1"),
				"f2":       New("em2", "ECode2"),
			},
		}
		assert.Equal(t, want, err)
	})
}

func Test_GenericFields_Filter(t *testing.T) {
	t.Run("filter", func(t *testing.T) {
		// --- Given ---
		key0 := errors.New("error")
		key2 := New("em2", "ECode")

		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"key0": key0,
				"key1": nil,
				"key2": key2,
			},
		}

		// --- When ---
		err := fs.Filter()

		// --- Then ---
		var fields *GenericFields[EDGeneric]
		assert.ErrorAs(t, &fields, err)
		assert.Len(t, 2, fields.fields)
		assert.Same(t, fields.fields["key0"], key0)
		assert.Same(t, fields.fields["key2"], key2)
	})

	t.Run("empty", func(t *testing.T) {
		// --- Given ---
		var fs GenericFields[EDGeneric]

		// --- When ---
		err := fs.Filter()

		// --- Then ---
		assert.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("nil instance", func(t *testing.T) {
		// --- Given ---
		var fs *GenericFields[EDGeneric]

		// --- When ---
		err := fs.Filter()

		// --- Then ---
		assert.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("all fields nil", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": nil,
				"f1": nil,
				"f2": nil,
			},
		}

		// --- When ---
		err := fs.Filter()

		// --- Then ---
		assert.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"key0": errors.New("error"),
				"key1": nil,
			},
		}

		// --- When ---
		_ = fs.Filter()

		// --- Then ---
		assert.Len(t, 2, fs.fields)
		_, _ = assert.HasKey(t, "key1", fs.fields)
	})

	t.Run("nested", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": nil,
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": nil,
							},
						},
					},
				},
				"f1": New("em2", "ECode2"),
				"f2": nil,
			},
		}

		// --- When ---
		err := fs.Filter()

		// --- Then ---
		var fe *GenericFields[EDGeneric]
		assert.ErrorAs(t, &fe, err)
		assert.Len(t, 2, fe.fields)
		_, _ = assert.HasKey(t, "f0", fe.fields)
		_, _ = assert.HasKey(t, "f1", fe.fields)

		f0, _ := assert.SameType(t, &GenericFields[EDGeneric]{}, fe.fields["f0"])
		assert.Len(t, 1, f0.fields)
		_, _ = assert.HasKey(t, "s0", f0.fields)

		have := must.Value(json.Marshal(err))
		want := `{
			"f0.s0":{"code":"ECode00","error":"em00"},
			"f1":{"code":"ECode2","error":"em2"}
		}`
		assert.JSON(t, want, string(have))
	})
}

func Test_GenericFields_Merge(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// --- Given ---
		f0 := map[string]error{
			"f0": errors.New("f0"),
		}
		f1 := map[string]error{
			"f1": errors.New("f1"),
		}
		fe := &GenericFields[EDGeneric]{fields: f0}

		// --- When ---
		fe.Merge(f1)

		// --- Then ---
		assert.Len(t, 2, f0)
		assert.Equal(t, "f0", f0["f0"].Error())
		assert.Equal(t, "f1", f0["f1"].Error())
	})

	t.Run("override", func(t *testing.T) {
		// --- Given ---
		f0 := map[string]error{
			"f0": nil,
		}
		f1 := map[string]error{
			"f0": errors.New("override"),
		}
		fe := &GenericFields[EDGeneric]{fields: f0}

		// --- When ---
		fe.Merge(f1)

		// --- Then ---
		assert.Len(t, 1, f0)
		assert.Equal(t, "override", f0["f0"].Error())
	})

	t.Run("no override", func(t *testing.T) {
		// --- Given ---
		f0 := map[string]error{
			"f0": errors.New("f0"),
		}
		f1 := map[string]error{
			"f0": errors.New("override"),
		}
		fe := &GenericFields[EDGeneric]{fields: f0}

		// --- When ---
		fe.Merge(f1)

		// --- Then ---
		assert.Len(t, 1, f0)
		assert.ErrorEqual(t, "f0", f0["f0"])
	})

	t.Run("nil receiver is no-op", func(t *testing.T) {
		// --- Given ---
		var fs *GenericFields[EDGeneric]

		// --- When --- Then --- (must not panic)
		fs.Merge(map[string]error{"f1": errors.New("msg1")})
	})

	t.Run("nil receiver with nil errs is no-op", func(t *testing.T) {
		// --- Given ---
		var fs *GenericFields[EDGeneric]

		// --- When --- Then --- (must not panic)
		fs.Merge(nil)
	})
}

func Fuzz_Fields_Get(f *testing.F) {
	// --- Given ---
	fs := &GenericFields[EDGeneric]{
		fields: map[string]error{
			"f0": &GenericFields[EDGeneric]{
				fields: map[string]error{
					"s0": New("f0.s0", "ECode00"),
					"s1": New("f0.s1", "ECode01"),
					"s2": &GenericFields[EDGeneric]{
						fields: map[string]error{
							"s0":            errors.New("f0.s2.s0"),
							"s1":            errors.New("f0.s2.s1"),
							"tag.name.mane": errors.New("f0.s2.tag.name.name"),
						},
					},
				},
			},
			"f1":       New("f1", "ECode1"),
			"f2.s0":    New("f2.s0", "ECode2"),
			"f2.s0.s0": New("f2.s0.s0", "ECode2"),
		},
	}

	tt := []string{
		"f0",
		"f1",
		"f0.s0",
		"f0.s2",
		"f0.s2.s0",
		"f2.s0.s0.s0",
		"f0.s2.tag.name.mane",
	}
	for _, tc := range tt {
		f.Add(tc)
	}

	// --- Then ---
	f.Fuzz(func(t *testing.T, s string) {
		_ = fs.Get(s)
	})
}

func Test_GenericFields_Get(t *testing.T) {
	// --- Given ---
	fs := &GenericFields[EDGeneric]{
		fields: map[string]error{
			"f0": &GenericFields[EDGeneric]{
				fields: map[string]error{
					"s0": New("f0.s0", "ECode00"),
					"s1": New("f0.s1", "ECode01"),
					"s2": &GenericFields[EDGeneric]{
						fields: map[string]error{
							"s0":            errors.New("f0.s2.s0"),
							"s1":            errors.New("f0.s2.s1"),
							"tag.name.mane": errors.New("f0.s2.tag.name.name"),
						},
					},
				},
			},
			"f1":       New("f1", "ECode1"),
			"f2.s0":    New("f2.s0", "ECode2"),
			"f2.s0.s0": New("f2.s0.s0", "ECode2"),
		},
	}

	t.Run("not nested", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f1")

		// --- Then ---
		assert.Equal(t, "f1", have.Error())
	})

	t.Run("nested", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f0")

		// --- Then ---
		want := "s0: f0.s0; " +
			"s1: f0.s1; " +
			"s2.s0: f0.s2.s0; " +
			"s2.s1: f0.s2.s1; " +
			"s2.tag.name.mane: f0.s2.tag.name.name"
		assert.Equal(t, want, have.Error())
	})

	t.Run("nested N1", func(t *testing.T) {
		// --- When ---
		err := fs.Get("f0.s0")

		// --- Then ---
		want := "f0.s0"
		assert.Equal(t, want, err.Error())
	})

	t.Run("nested N1 with sub", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f0.s2")

		// --- Then ---
		want := "s0: f0.s2.s0; " +
			"s1: f0.s2.s1; " +
			"tag.name.mane: f0.s2.tag.name.name"
		assert.Equal(t, want, have.Error())
	})

	t.Run("nested N2", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f0.s2.s0")

		// --- Then ---
		want := "f0.s2.s0"
		assert.Equal(t, want, have.Error())
	})

	t.Run("flatten case 1", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f2.s0")

		// --- Then ---
		want := "f2.s0"
		assert.Equal(t, want, have.Error())
	})

	t.Run("flatten case 2", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f2.s0.s0")

		// --- Then ---
		want := "f2.s0.s0"
		assert.Equal(t, want, have.Error())
	})

	t.Run("nested not existing case 1", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f0.s3")

		// --- Then ---
		assert.NoError(t, have)
	})

	t.Run("nested not existing case 2", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f2.s0.s0.s0")

		// --- Then ---
		assert.NoError(t, have)
	})

	t.Run("tag name", func(t *testing.T) {
		// --- When ---
		have := fs.Get("f0.s2.tag.name.mane")

		// --- Then ---
		want := "f0.s2.tag.name.name"
		assert.Equal(t, want, have.Error())
	})

	t.Run("the key is not a dot-path prefix of a field", func(t *testing.T) {
		// --- Given ---
		fs2 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"a": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"c": errors.New("nested"),
					},
				},
			},
		}

		// --- When ---
		have := fs2.Get("abc")

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_GenericFields_Set(t *testing.T) {
	t.Run("set new field", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{fields: map[string]error{}}

		// --- When ---
		fs.Set("f0", ErrTst)

		// --- Then ---
		assert.Same(t, ErrTst, fs.fields["f0"])
	})

	t.Run("nil fields map is initialized", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{}

		// --- When ---
		fs.Set("f0", ErrTst)

		// --- Then ---
		assert.Same(t, ErrTst, fs.fields["f0"])
	})

	t.Run("overwrite existing field", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{fields: map[string]error{"f0": errors.New("old")}}

		// --- When ---
		fs.Set("f0", ErrTst)

		// --- Then ---
		assert.Same(t, ErrTst, fs.fields["f0"])
	})

	t.Run("set nil error", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{fields: map[string]error{"f0": ErrTst}}

		// --- When ---
		fs.Set("f0", nil)

		// --- Then ---
		assert.Nil(t, fs.fields["f0"])
	})

	t.Run("nil receiver is no-op", func(t *testing.T) {
		// --- Given ---
		var fs *GenericFields[EDGeneric]

		// --- When --- Then --- (must not panic)
		fs.Set("f0", ErrTst)
	})
}

func Test_GenericFields_Len(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{fields: map[string]error{
			"f0": errors.New("em0"),
			"f1": errors.New("em1"),
		}}

		// --- When ---
		have := fs.Len()

		// --- Then ---
		assert.Equal(t, 2, have)
	})

	t.Run("empty fields", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{fields: map[string]error{}}

		// --- When ---
		have := fs.Len()

		// --- Then ---
		assert.Equal(t, 0, have)
	})

	t.Run("nil receiver returns zero", func(t *testing.T) {
		// --- Given ---
		var fs *GenericFields[EDGeneric]

		// --- When ---
		have := fs.Len()

		// --- Then ---
		assert.Equal(t, 0, have)
	})
}

func Test_GenericFields_MarshalJSON(t *testing.T) {
	t.Run("with many levels", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": New("em00", "ECode00"),
						"s1": New("em01", "ECode01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		have, err := fs.MarshalJSON()

		// --- Then ---
		assert.NoError(t, err)
		want := `{
			"f0.s0": {"code": "ECode00", "error": "em00"},
			"f0.s1": {"code": "ECode01", "error": "em01"},
			"f0.s2.s0": {"code": "ECGeneric", "error": "em020"},
			"f1": {"code": "ECode1", "error": "em1", "meta": {"key": "val"}},
			"f2": {"code": "ECode2", "error": "em2"}
		}`
		assert.JSON(t, want, string(have))
	})

	t.Run("all fields nil returns empty object", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": nil,
				"f1": nil,
			},
		}

		// --- When ---
		have, err := fs.MarshalJSON()

		// --- Then ---
		assert.NoError(t, err)
		assert.JSON(t, `{}`, string(have))
	})

	t.Run("field marshaller error", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &TErrMarshalJSON{errors.New("abc")},
			},
		}

		// --- When ---
		have, err := fs.MarshalJSON()

		// --- Then ---
		assert.ErrorEqual(t, "abc", err)
		assert.Nil(t, have)
	})
}

func Test_GenericFields_UnmarshalJSON(t *testing.T) {
	t.Run("round-trip", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": New("em0", "ECode0"),
				"f1": New("em1", "ECode1", Meta().Str("key", "val").Option()),
			},
		}
		data, err := fs.MarshalJSON()
		assert.NoError(t, err)

		// --- When ---
		var got GenericFields[EDGeneric]
		err = json.Unmarshal(data, &got)

		// --- Then ---
		assert.NoError(t, err)
		assert.Len(t, 2, got.fields)
		assert.ErrorEqual(t, "em0", got.fields["f0"])
		assert.ErrorEqual(t, "em1", got.fields["f1"])
	})

	t.Run("empty object yields empty fields", func(t *testing.T) {
		// --- Given ---
		var got GenericFields[EDGeneric]

		// --- When ---
		err := json.Unmarshal([]byte(`{}`), &got)

		// --- Then ---
		assert.NoError(t, err)
		assert.Len(t, 0, got.fields)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		// --- Given ---
		var got GenericFields[EDGeneric]

		// --- When ---
		err := json.Unmarshal([]byte(`not-json`), &got)

		// --- Then ---
		assert.Error(t, err)
	})
}

func Test_Flatten(t *testing.T) {
	t.Run("flatten single", func(t *testing.T) {
		// --- Given ---
		fs := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": errors.New("em00"),
						"s1": errors.New("em01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1"),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		err := Flatten[EDGeneric](fs)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0.s0":    errors.New("em00"),
				"f0.s1":    errors.New("em01"),
				"f0.s2.s0": errors.New("em020"),
				"f1":       New("em1", "ECode1"),
				"f2":       New("em2", "ECode2"),
			},
		}
		assert.Equal(t, want, err)
	})

	t.Run("flatten multiple", func(t *testing.T) {
		// --- Given ---
		fs0 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0": &GenericFields[EDGeneric]{
					fields: map[string]error{
						"s0": errors.New("em00"),
						"s1": errors.New("em01"),
						"s2": &GenericFields[EDGeneric]{
							fields: map[string]error{
								"s0": errors.New("em020"),
							},
						},
					},
				},
				"f1": New("em1", "ECode1"),
			},
		}

		fs1 := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f1": New("other", "ECOther"),
				"f2": New("em2", "ECode2"),
			},
		}

		// --- When ---
		err := Flatten[EDGeneric](fs0, fs1)

		// --- Then ---
		want := &GenericFields[EDGeneric]{
			fields: map[string]error{
				"f0.s0":    errors.New("em00"),
				"f0.s1":    errors.New("em01"),
				"f0.s2.s0": errors.New("em020"),
				"f1":       New("other", "ECOther"),
				"f2":       New("em2", "ECode2"),
			},
		}
		assert.Equal(t, want, err)
	})
}
