package reflect

import (
	"fmt"
	"reflect"
)

func newCannotCastErr(from, to reflect.Type) error {
	return fmt.Errorf("cannot cast `%s` to `%s`", from.String(), to.String())
}

// Convert converts given value to given type whenever it is possible.
// In opposition to built-in reflect package it allows to convert []interface{} to []type.
// todo add comment about zero values
// todo []interface{} => []type is possible, add []type => []interface{}
func Convert(value interface{}, to reflect.Type) (reflect.Value, error) {
	// it is required to avoid panic (reflect: call of reflect.Value.Type on zero Value)
	// in case of the following code
	// caller.MustCall(func(v interface{}) { fmt.Println(v) }, v)
	if value == nil {
		return reflect.Zero(to), nil
	}
	from := reflect.ValueOf(value)
	if from.Kind() == reflect.Invalid {
		return reflect.Value{}, newCannotCastErr(from.Type(), to)
	}
	if from.Type().ConvertibleTo(to) {
		return from.Convert(to), nil
	}

	slice, ok := convertSliceInterface(from, to)
	if !ok {
		return reflect.Value{}, newCannotCastErr(from.Type(), to)
	}

	return slice, nil
}

func convertSliceInterface(from reflect.Value, to reflect.Type) (reflect.Value, bool) {
	if from.Kind() != reflect.Slice || from.Type().Elem().Kind() != reflect.Interface || to.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}
	cp := reflect.MakeSlice(to, 0, 0)
	for i := 0; i < from.Len(); i++ {
		item := from.Index(i)
		for item.Kind() == reflect.Interface {
			item = item.Elem()
		}
		if !item.Type().ConvertibleTo(to.Elem()) {
			return reflect.Value{}, false
		}
		cp = reflect.Append(cp, item.Convert(to.Elem()))
	}
	return cp, true
}
