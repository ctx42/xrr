package xrr

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Fields represents a collection of errors that are indexed by field names.
//
// The field (name) indexed errors are mostly useful for validation errors.
type Fields map[string]error

// FieldError return a new instance of [Fields] with the given error. Returns
// nil when the error is nil.
func FieldError(field string, err error) error {
	if err == nil {
		return nil
	}
	return Fields{field: err}
}

// AddField sets new error on provided [Fields] instance. If err is nil, the
// call is no-op. If ers is nil and err is not, the new instance of [Fields]
// is created and the passed error is added to it with a given field name.
func AddField(ers *Fields, field string, err error) {
	if err == nil {
		return
	}
	if *ers == nil {
		*ers = Fields{}
	}
	(*ers)[field] = err
}

// GetFields returns field errors if the error implements the [Fielder]
// interface. Otherwise, it returns nil.
func GetFields(err error) map[string]error {
	if mg, ok := err.(Fielder); ok {
		return mg.ErrFields()
	}
	return nil
}

// GetFieldError returns error for given field name. It expects the error to
// be an instance of [Fields]. Returns nil when err is nil, not an instance of
// [Fields] or when there is no error for the given field name.
func GetFieldError(err error, field string) error {
	if fs := GetFields(err); fs != nil {
		return get(fs, field)
	}
	return nil
}

// FieldErrorIs returns true if err is an instance of [Fields] with the given
// field name and the [errors.Is] returns true for the error and the target.
func FieldErrorIs(err error, field string, target error) bool {
	return errors.Is(GetFieldError(err, field), target)
}

// FieldNames returns alphabetically sorted fields names if the error is an
// instance of [Fields]. Otherwise, it returns nil.
func FieldNames(err error) []string {
	fs := GetFields(err)
	if fs == nil {
		return nil
	}
	var names []string
	for name, _ := range fs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// FieldRename if the error an instance of [Fields], it renames the given field
// name from "from" to "to". If err is not a [Fields] instance, it is no-op.
func FieldRename(err error, from, to string) {
	var fe Fields
	if from != to && errors.As(err, &fe) {
		if e, ok := fe[from]; ok {
			fe[to] = e
			delete(fe, from)
		}
	}
}

// MergeFields merges multiple instances of [Fields]. Expects all errors to be
// instances [Fields]. But will handle other error types as well by creating
// fake indexed fields for them. The nil instances os [Fields] are ignored, but
// nil field errors are kept. The field errors with the same name are
// overwritten by errors that are later in the list.
func MergeFields(ers ...error) error {
	if err := mergeFields(ers...); err != nil {
		return err
	}
	return nil
}
func mergeFields(ers ...error) Fields {
	if len(ers) == 0 {
		return nil
	}
	var i int
	var k string
	var e error
	var first, fe Fields

	// Find first non-nil error.
	for i, e = range ers {
		if e == nil {
			continue
		}
		if !errors.As(e, &first) {
			key := fmt.Sprintf("__field__%d", i)
			first = Fields{key: e}
		}
		break
	}

	// All errors were nil.
	if first == nil {
		return nil
	}

	for j, err := range ers[i:] {
		if err == nil {
			continue
		}
		if !errors.As(err, &fe) {
			key := fmt.Sprintf("__field__%d", i+j)
			first[key] = err
		}
		for k, e = range fe {
			// Don't overwrite existing non-nil field errors with nil errors.
			if existing := first[k]; existing != nil && e == nil {
				continue
			}
			first[k] = e
		}
	}

	return first
}

func (fs Fields) ErrFields() map[string]error { return fs }

// Error returns string representation of field errors.
func (fs Fields) Error() string {
	return formatFields(fs.ErrFields(), false)
}

// Is implements the interface used by [errors.Is].
func (fs Fields) Is(other error) bool {
	if other == nil {
		return false
	}
	for _, e := range fs {
		if errors.Is(e, other) {
			return true
		}
	}
	return false
}

func (fs Fields) Format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'q':
		msg := fs.Error()
		if verb == 'q' {
			msg = fmt.Sprintf("%q", msg)
		}
		_, _ = fmt.Fprint(state, msg)

	case 'v':
		if state.Flag('+') {
			_, _ = fmt.Fprint(state, formatFields(fs.ErrFields(), true))
		} else {
			msg := fs.Error()
			_, _ = fmt.Fprint(state, msg)
		}
	}
}

// Flatten flattens a nested map of errors to single one level map. The fields
// for nested errors are prefixed with the field name of the parent separated
// by dots (.).
//
// Flatten example:
//
//	map[string]error{
//	  "a": errors.New("a"),
//	  "a.b": errors.New("b"),
//	}
func (fs Fields) Flatten() Fields {
	visitor := make(map[string]error, len(fs))
	flatten(visitor, "", fs)
	return visitor
}

// Filter removes all keys with nil values from Fields and returns it as an
// error. If the length of Fields becomes 0, it will return nil.
func (fs Fields) Filter() error {
	if fs == nil {
		return nil
	}
	if ret := filterMap(fs); ret != nil {
		return Fields(ret)
	}
	return nil
}

// filterMap removes all keys with nil values and returns it. If the length of
// the filtered map is zero, it will return nil.
func filterMap(fs map[string]error) map[string]error {
	for key, value := range fs {
		if value == nil {
			delete(fs, key)
			continue
		}
		var fields Fields
		if errors.As(value, &fields) {
			if filtered := filterMap(fields); filtered != nil {
				fs[key] = Fields(filtered)
			} else {
				delete(fs, key)
			}
		}
	}
	if len(fs) == 0 {
		return nil
	}
	return fs
}

// Merge adds all non nil errors from errs overriding already existing keys.
func (fs Fields) Merge(errs map[string]error) Fields {
	if fs == nil && len(errs) == 0 {
		return nil
	}
	if fs == nil {
		//goland:noinspection GoAssignmentToReceiver
		fs = Fields{}
	}
	for key, err := range errs {
		if fs[key] == nil {
			fs[key] = err
		}
	}
	return fs
}

// Get returns an error for the given field, nil if the field does not exist.
func (fs Fields) Get(field string) error {
	return get(fs, field)
}

// get returns an error for the given field, nil if the field does not exist.
func get(ers map[string]error, field string) error {
	sub := field
	for key, err := range ers {
		if sub == key {
			return err
		}
		if !strings.HasPrefix(sub, key) {
			continue
		}
		var fs Fields
		if errors.As(err, &fs) {
			sub = sub[len(key)+1:]
			if err := get(fs, sub); err != nil {
				return err
			}
		}
	}
	return nil
}

func (fs Fields) MarshalJSON() ([]byte, error) {
	visitor := make(map[string]error, len(fs))
	flatten(visitor, "", fs)
	ret := make(map[string]json.RawMessage, len(visitor))
	for k, v := range filterMap(visitor) {
		if e, ok := v.(json.Marshaler); ok { // nolint: errorlint
			data, err := e.MarshalJSON()
			if err != nil {
				return nil, err
			}
			ret[k] = data
			continue
		}
		data, err := json.Marshal(Wrap(v))
		if err != nil {
			return nil, err
		}
		ret[k] = data
	}

	return json.Marshal(ret)
}

// Flatten first merges all the provided errors, then it flattens a nested map
// of errors to single one level map. The fields for nested errors are prefixed
// with the field name of the parent separated by dots (.).
//
// Flatten example:
//
//	map[string]error{
//	  "a": errors.New("a"),
//	  "a.b": errors.New("b"),
//	}
func Flatten(err ...error) error {
	visitor := make(map[string]error)
	fls := mergeFields(err...)
	flatten(visitor, "", fls)
	return Fields(visitor)
}

// flatten flattens nested map of errors.
func flatten(visitor map[string]error, pref string, fields map[string]error) {
	for field, err := range fields {
		var fs Fields
		if errors.As(err, &fs) {
			flatten(visitor, prefix(pref, field), fs.ErrFields())
			continue
		}
		visitor[prefix(pref, field)] = err
	}
}

// formatFields returns string representation of Fields.
func formatFields(fs map[string]error, codes bool) string {
	if len(fs) == 0 {
		return ""
	}

	visitor := make(map[string]error, len(fs))
	flatten(visitor, "", fs)

	keys := make([]string, len(visitor))
	i := 0
	for key := range visitor {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	var s strings.Builder
	for i, key := range keys {
		err := visitor[key]
		if err == nil {
			continue
		}
		if i > 0 {
			s.WriteString("; ")
		}
		if codes {
			_, _ = fmt.Fprintf(
				&s,
				"%v: %v (%s)",
				key,
				err.Error(),
				GetCode(err),
			)
		} else {
			_, _ = fmt.Fprintf(
				&s,
				"%v: %v",
				key,
				err.Error(),
			)
		}
	}
	return s.String()
}
