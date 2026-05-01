// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package xrr

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Compile time checks.
var (
	_ error            = (*GenericFields[EDXrr])(nil)
	_ Fielder          = (*GenericFields[EDXrr])(nil)
	_ json.Marshaler   = (*GenericFields[EDXrr])(nil)
	_ json.Unmarshaler = (*GenericFields[EDXrr])(nil)
)

// GenericFields represents a generic type for creating domain-specific
// field-indexed errors.
type GenericFields[T Domain] struct {
	fields map[string]error
}

// NewFields creates a new [GenericFields][T] from the given map. The map
// is stored directly without copying.
func NewFields[T Domain](fields map[string]error) *GenericFields[T] {
	return &GenericFields[T]{fields: fields}
}

// FieldsFunc returns a function for creating domain-specific field errors.
func FieldsFunc[T Domain]() func(field string, err error) *GenericFields[T] {
	return func(field string, err error) *GenericFields[T] {
		if err == nil {
			return nil
		}
		return &GenericFields[T]{fields: map[string]error{field: err}}
	}
}

// GetFields returns field errors if the error implements the [Fielder]
// interface. Otherwise, it returns nil.
func GetFields(err error) map[string]error {
	if mg, ok := err.(Fielder); ok {
		return mg.ErrorFields()
	}
	return nil
}

// GetFieldError returns an error for the given field name. It expects the
// error to implement [Fielder]. Returns nil when err is nil, does not
// implement [Fielder], or has no error for the given field name.
func GetFieldError(err error, field string) error {
	if fs := GetFields(err); fs != nil {
		return get(fs, field)
	}
	return nil
}

// FieldErrorIs returns true if err implements [Fielder], has the given field
// name, and [errors.Is] returns true for that field's error and target.
func FieldErrorIs(err error, field string, target error) bool {
	return errors.Is(GetFieldError(err, field), target)
}

// FieldNames returns alphabetically sorted field names if the error implements
// [Fielder]. Otherwise, it returns nil.
func FieldNames(err error) []string {
	fs := GetFields(err)
	if fs == nil {
		return nil
	}
	var names []string
	for name := range fs {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// MergeFields merges multiple [Fielder] errors into a single [GenericFields[T]].
// Returns nil if all inputs are nil or the slice is empty.
//
// Merge rules:
//   - nil errors are skipped.
//   - For [Fielder] errors, all field entries are merged into the result.
//   - For non-[Fielder] errors, a synthetic key "__field__N" is used, where N
//     is the position of the error in the argument list.
//   - When two inputs share a field name, the later one wins.
//   - nil field values are preserved (they are not treated as absent).
func MergeFields[T Domain](ers ...error) error {
	if fe := mergeFields(ers...); fe != nil {
		return &GenericFields[T]{fields: fe}
	}
	return nil
}

// mergeFields merges multiple field error maps into a single map. Non-[Fielder]
// errors are assigned synthetic keys of the form "__field__N".
func mergeFields(ers ...error) map[string]error {
	if len(ers) == 0 {
		return nil
	}

	var i int
	var k string
	var e error
	var ok bool
	var fe Fielder
	var first map[string]error

	// Find first non-nil error.
	for i, e = range ers {
		if e == nil {
			continue
		}
		if fe, ok = e.(Fielder); ok {
			first = fe.ErrorFields()
		} else {
			key := fmt.Sprintf("__field__%d", i)
			first = map[string]error{key: e}
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
		if fe, ok = err.(Fielder); ok {
			for k, e = range fe.ErrorFields() {
				// Don't overwrite existing non-nil field with nil errors.
				if existing := first[k]; existing != nil && e == nil {
					continue
				}
				first[k] = e
			}
		} else {
			key := fmt.Sprintf("__field__%d", i+j)
			first[key] = err
		}
	}

	return first
}

func (fs *GenericFields[T]) ErrorFields() map[string]error { return fs.fields }

func (fs *GenericFields[T]) Error() string {
	return formatFields(fs.ErrorFields(), false)
}

func (fs *GenericFields[T]) Unwrap() []error {
	flat := fs.Flatten()
	fields, ers := sortFields(flat.fields)
	var j int
	for i, err := range ers {
		if err != nil {
			ers[j] = fmt.Errorf("%s: %w", fields[i], err)
			j++
		}
	}
	return ers[:j]
}

// Is implements the interface used by [errors.Is].
func (fs *GenericFields[T]) Is(other error) bool {
	if other == nil {
		return false
	}
	for _, e := range fs.fields {
		if errors.Is(e, other) {
			return true
		}
	}
	return false
}

func (fs *GenericFields[T]) Format(state fmt.State, verb rune) {
	switch verb {
	case 's', 'q':
		msg := fs.Error()
		if verb == 'q' {
			msg = fmt.Sprintf("%q", msg)
		}
		_, _ = fmt.Fprint(state, msg)

	case 'v':
		if state.Flag('+') {
			_, _ = fmt.Fprint(state, formatFields(fs.ErrorFields(), true))
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
func (fs *GenericFields[T]) Flatten() *GenericFields[T] {
	visitor := make(map[string]error, len(fs.fields))
	flatten(visitor, "", fs.fields)
	return &GenericFields[T]{fields: visitor}
}

// Filter removes all keys with nil values from Fields and returns it as an
// error. If the length of Fields becomes 0, it will return nil.
func (fs *GenericFields[T]) Filter() error {
	if fs == nil {
		return nil
	}
	if ret := filterMap[T](fs.fields); ret != nil {
		return ret
	}
	return nil
}

// filterMap returns a new map with nil values removed. Nested [Fielder] values
// are filtered recursively. Returns nil if no entries survive filtering.
func filterMap[T Domain](fs map[string]error) *GenericFields[T] {
	ret := make(map[string]error, len(fs))
	for key, value := range fs {
		if value == nil {
			continue
		}
		if fls, ok := value.(Fielder); ok {
			if filtered := filterMap[T](fls.ErrorFields()); filtered != nil {
				ret[key] = filtered
			}
			continue
		}
		ret[key] = value
	}
	if len(ret) == 0 {
		return nil
	}
	return &GenericFields[T]{fields: ret}
}

// Merge adds errors from errs for keys that are not already set in fs.
// It is a no-op when fs is nil or errs is empty.
func (fs *GenericFields[T]) Merge(errs map[string]error) {
	if fs == nil || len(errs) == 0 {
		return
	}
	if fs.fields == nil {
		fs.fields = make(map[string]error, len(errs))
	}
	for key, err := range errs {
		if fs.fields[key] == nil {
			fs.fields[key] = err
		}
	}
}

// Get returns an error for the given field, nil if the field does not exist.
func (fs *GenericFields[T]) Get(field string) error {
	return get(fs.fields, field)
}

// get returns an error for the given field, nil if the field does not exist.
func get(ers map[string]error, field string) error {
	for key, err := range ers {
		if field == key {
			return err
		}
		suffix, ok := strings.CutPrefix(field, key+".")
		if !ok {
			continue
		}
		if fls, ok := err.(Fielder); ok {
			if err = get(fls.ErrorFields(), suffix); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set sets the error for the given field. It is a no-op when fs is nil.
func (fs *GenericFields[T]) Set(field string, err error) {
	if fs == nil {
		return
	}
	if fs.fields == nil {
		fs.fields = make(map[string]error)
	}
	fs.fields[field] = err
}

// Len returns the number of fields. Returns 0 if fs is nil.
func (fs *GenericFields[T]) Len() int {
	if fs == nil {
		return 0
	}
	return len(fs.fields)
}

func (fs *GenericFields[T]) MarshalJSON() ([]byte, error) {
	visitor := make(map[string]error, len(fs.fields))
	flatten(visitor, "", fs.fields)
	ret := make(map[string]json.RawMessage, len(visitor))
	fls := filterMap[T](visitor)
	if fls == nil {
		return []byte(`{}`), nil
	}
	for k, v := range fls.fields {
		if e, ok := v.(json.Marshaler); ok { // nolint: errorlint
			data, err := e.MarshalJSON()
			if err != nil {
				return nil, err
			}
			ret[k] = data
			continue
		}
		m := errorAsMap(v)
		data, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		ret[k] = data
	}
	return json.Marshal(ret)
}

func (fs *GenericFields[T]) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	fs.fields = make(map[string]error, len(raw))
	for k, v := range raw {
		var e GenericError[T]
		if err := json.Unmarshal(v, &e); err != nil {
			return err
		}
		fs.fields[k] = &e
	}
	return nil
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
func Flatten[T Domain](err ...error) error {
	visitor := make(map[string]error)
	fls := mergeFields(err...)
	flatten(visitor, "", fls)
	return &GenericFields[T]{fields: visitor}
}

// flatten flattens nested map of errors.
func flatten(visitor map[string]error, pref string, fields map[string]error) {
	for field, err := range fields {
		if fls, ok := err.(Fielder); ok {
			flatten(visitor, prefix(pref, field), fls.ErrorFields())
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
	slices.Sort(keys)

	var s strings.Builder
	for _, key := range keys {
		err := visitor[key]
		if err == nil {
			continue
		}
		if s.Len() > 0 {
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
				errorMessage(err),
			)
		}
	}
	return s.String()
}
