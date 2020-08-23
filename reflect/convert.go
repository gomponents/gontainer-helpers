package reflect

import (
	"reflect"
)

func Convert(from reflect.Value, to reflect.Type) (reflect.Value, bool) {
	if from.Type().ConvertibleTo(to) {
		return from.Convert(to), true
	}

	return convertSliceInterface(from, to)
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
