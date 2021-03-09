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

	v, err := ref.Convert(val, f.Type())
	if err != nil {
		// todo wrap err
		return err
	}

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

	wrap := func(err error) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("set `%T`.`%s`: %s", strct, field, err.Error())
	}

	// chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Ptr (, reflect.Ptr...), reflect.Struct)
	isInterfaceOverPointerChain := func(chain kindChain) bool {
		if len(chain) < 4 {
			return false
		}
		if chain[0] != reflect.Ptr {
			return false
		}
		if chain[1] != reflect.Interface {
			return false
		}
		if chain[len(chain)-1] != reflect.Struct {
			return false
		}

		for _, c := range chain[2 : len(chain)-2] {
			if c != reflect.Ptr {
				return false
			}
		}

		return true
	}

	switch {
	// s := struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Struct):
		return wrap(set(
			reflectVal.Elem(),
			field,
			val,
		))

	// case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Ptr (, reflect.Ptr...), reflect.Struct):
	// var s interface{} = &struct{ val int }{}
	// Set(&s...
	case isInterfaceOverPointerChain(chain):
		elem := reflectVal.Elem()
		for i := 0; i < len(chain)-2; i++ {
			elem = elem.Elem()
		}
		return wrap(set(elem, field, val))

	// var s interface{} = struct{ val int }{}
	// Set(&s...
	case chain.equalTo(reflect.Ptr, reflect.Interface, reflect.Struct):
		v := reflectVal.Elem()
		tmp := reflect.New(v.Elem().Type()).Elem()
		tmp.Set(v.Elem())
		if err := wrap(set(tmp, field, val)); err != nil {
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

func valueToKindChain(val interface{}) kindChain {
	var r kindChain
	v := reflect.ValueOf(val)
	for {
		r = append(r, v.Kind())
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
			continue
		}
		break
	}
	return r
}
