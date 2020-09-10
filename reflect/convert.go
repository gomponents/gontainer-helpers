package reflect

import (
	"fmt"
	"reflect"
)

func newCannotCastErr(from, to reflect.Type) error {
	return fmt.Errorf("cannot cast `%s` to `%s`", from.String(), to.String())
}

func isNilable(k reflect.Kind) bool {
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	}
	return false
}

// Convert converts given value to given type whenever it is possible.
// In opposition to built-in reflect package it allows to convert []interface{} to []type and []type to []interface{}.
func Convert(value interface{}, to reflect.Type) (reflect.Value, error) {
	// it is required to avoid panic (reflect: call of reflect.Value.Type on zero Value)
	// in case of the following code
	// caller.MustCall(func(v interface{}) { fmt.Println(v) }, v)
	if value == nil {
		if isNilable(to.Kind()) {
			return reflect.Zero(to), nil
		}
		return reflect.Value{}, fmt.Errorf("cannot cast `%T` to `%s`", value, to.String())
	}
	from := reflect.ValueOf(value)
	if from.Type().ConvertibleTo(to) {
		return from.Convert(to), nil
	}

	slice, ok := convertSlice(from, to)
	if !ok {
		return reflect.Value{}, newCannotCastErr(from.Type(), to)
	}

	return slice, nil
}

func convertSlice(from reflect.Value, to reflect.Type) (reflect.Value, bool) {
	if from.Kind() != reflect.Slice || to.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}

	canCastSlice := func() bool {
		if from.Type().Elem().Kind() == reflect.Interface || to.Elem().Kind() == reflect.Interface {
			return true
		}

		if from.Type().Elem().ConvertibleTo(to.Elem()) {
			return true
		}

		return false
	}

	if !canCastSlice() {
		return reflect.Value{}, false
	}

	cp := reflect.MakeSlice(to, 0, 0)
	for i := 0; i < from.Len(); i++ {
		item := from.Index(i)
		for item.Kind() == reflect.Interface {
			item = item.Elem()
		}
		if !item.IsValid() { // nil
			if isNilable(to.Elem().Kind()) {
				cp = reflect.Append(cp, reflect.Zero(to.Elem()))
				continue
			}
			return reflect.Value{}, false
		}
		if !item.Type().ConvertibleTo(to.Elem()) {
			return reflect.Value{}, false
		}
		cp = reflect.Append(cp, item.Convert(to.Elem()))
	}
	return cp, true
}
