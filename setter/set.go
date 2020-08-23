package setter

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	ref "github.com/gomponents/gontainer-helpers/reflect"
)

type kindChain []reflect.Kind

func set(strct reflect.Value, field string, val interface{}) error {
	f := strct.FieldByName(field)
	if !f.IsValid() {
		return fmt.Errorf("field `%s` does not exist", field)
	}
	var v reflect.Value
	if val == nil {
		v = reflect.Zero(f.Type()) // required for nil
	} else {
		v = reflect.ValueOf(val)
	}

	cp, ok := ref.Convert(v, f.Type())
	if !ok {
		return fmt.Errorf("cannot cast `%s` to `%s`", v.Type().String(), f.Type().String())
	}
	v = cp

	if !f.CanSet() { // handle unexported fields
		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	}
	f.Set(v)
	return nil
}

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
	chain := valueToKindChain(strct)

	// removes prepending duplicate Ptr elements
	// e.g.
	// s := &struct{ val int }{}
	// Set(&s... // chain == {Ptr, Ptr, Struct}
	reflectVal := reflect.ValueOf(strct)
	for len(chain) >= 2 && chain[0] == reflect.Ptr && chain[1] == reflect.Ptr {
		reflectVal = reflectVal.Elem()
		chain = chain[1:]
	}

	switch {
	// s := struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Struct):
		return set(
			reflectVal.Elem(),
			field,
			val,
		)

	// var s interface{} = &struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Ptr, reflect.Struct):
		return set(
			reflectVal.Elem().Elem().Elem(),
			field,
			val,
		)

	// var s interface{} = struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Struct):
		v := reflectVal.Elem()
		tmp := reflect.New(v.Elem().Type()).Elem()
		tmp.Set(v.Elem())
		if err := set(tmp, field, val); err != nil {
			return err
		}
		v.Set(tmp)
		return nil

	default:
		return fmt.Errorf("invalid parameter, setter.Set expects pointer to struct, given %s", chain.String())
	}
}

// MustSet calls Set with given parameters, it panics in case of error.
func MustSet(strct interface{}, field string, val interface{}) {
	err := Set(strct, field, val)
	if err != nil {
		panic(err)
	}
}

func (c kindChain) equalTo(kinds ...reflect.Kind) bool {
	if len(c) != len(kinds) {
		return false
	}

	for i := 0; i < len(c); i++ {
		if c[i] != kinds[i] {
			return false
		}
	}

	return true
}

func (c kindChain) String() string {
	parts := make([]string, len(c))
	for i, k := range c {
		parts[i] = k.String()
	}
	return strings.Join(parts, ".")
}

func valueToKindChain(v interface{}) kindChain {
	var r kindChain
	ref := reflect.ValueOf(v)
	for {
		r = append(r, ref.Kind())
		if ref.Kind() == reflect.Ptr || ref.Kind() == reflect.Interface {
			ref = ref.Elem()
			continue
		}
		break
	}
	return r
}
