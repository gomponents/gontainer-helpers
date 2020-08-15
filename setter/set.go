package setter

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Set assigns value of `val` to field `field` on struct `strct`.
// For `val` == nil it will always assign zero-value (e.g. 0 for int).
// Unexported fields are supported.
//
// Example of usage:
// a := struct {
//     val int
// }{}
// err := Set(&a, "val", 5)
// fmt.Println(err) // nil
// fmt.Println(a.val) // 5
func Set(strct interface{}, field string, val interface{}) error {
	s := reflect.ValueOf(strct)

	if s.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("expects `%s`, `%s` given", reflect.Ptr.String(), s.Type().Kind().String())
	}

	if s.Elem().Kind() == reflect.Interface {
		s = s.Elem()
	}

	if s.Elem().Kind() == reflect.Ptr {
		s = s.Elem()
	}

	if s.Type().Kind() != reflect.Ptr {
		return fmt.Errorf(
			"invalid arg (given `var arg interface{} := %s{}`, replace by `var arg interface{} := &%s{}`)",
			s.Elem().Type(),
			s.Elem().Type(),
		)
	}

	if s.Type().Elem().Kind() != reflect.Struct {
		return fmt.Errorf("invalid pointer dest: expects `%s`, `%s` given", reflect.Struct.String(), s.Type().Elem().Kind().String())
	}

	f := s.Elem().FieldByName(field)
	if !f.IsValid() {
		return fmt.Errorf("field `%s` does not exist", field)
	}

	var v reflect.Value
	if val == nil {
		v = reflect.Zero(f.Type()) // required for nil
	} else {
		v = reflect.ValueOf(val)
	}
	if !v.Type().ConvertibleTo(f.Type()) {
		return fmt.Errorf("cannot cast `%s` to `%s`", v.Type().Kind().String(), f.Type().Kind().String())
	}
	v = v.Convert(f.Type())

	if !f.CanSet() { // handle unexported fields
		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	}

	f.Set(v)
	return nil
}

// MustSet calls Set with given parameters, it panics in case of error.
func MustSet(strct interface{}, field string, val interface{}) {
	err := Set(strct, field, val)
	if err != nil {
		panic(err)
	}
}
