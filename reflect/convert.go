package reflect

import (
	"fmt"
	"reflect"
)

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

	if !isConvertibleSlice(from.Type(), to) {
		return reflect.Value{}, fmt.Errorf("cannot cast `%s` to `%s`", from.Type().String(), to.String())
	}

	slice, err := convertSlice(from, to)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("cannot cast `%s` to `%s`: %s", from.Type().String(), to.String(), err.Error())
	}

	return slice, nil
}

func isConvertibleSlice(from reflect.Type, to reflect.Type) bool {
	if from.Kind() != reflect.Slice || to.Kind() != reflect.Slice {
		return false
	}

	if from.Elem().Kind() == reflect.Interface || to.Elem().Kind() == reflect.Interface {
		return true
	}

	if from.Elem().ConvertibleTo(to.Elem()) {
		return true
	}

	if isConvertibleSlice(from.Elem(), to.Elem()) {
		return true
	}

	return false
}

func convertSlice(from reflect.Value, to reflect.Type) (reflect.Value, error) {
	cp := reflect.MakeSlice(to, 0, 0)
	for i := 0; i < from.Len(); i++ {
		item := from.Index(i)
		for item.Kind() == reflect.Interface {
			item = item.Elem()
		}
		var currVal interface{} = nil
		if item.IsValid() {
			currVal = item.Interface()
		}
		curr, err := Convert(currVal, to.Elem())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("el%d: %s", i, err.Error())
		}
		cp = reflect.Append(cp, curr)
	}
	return cp, nil
}
