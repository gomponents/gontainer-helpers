package reflect

import (
	"fmt"
	"reflect"
)

func MustCopy(from interface{}, to interface{}) {
	if err := Copy(from, to); err != nil {
		panic(err)
	}
}

// Copy sets value of `from` to `to`
// from := 5
// b := 0
// Copy(from, &to)
// fmt.Println(to) // 5
func Copy(from interface{}, to interface{}) error {
	f := reflect.ValueOf(from)
	t := reflect.ValueOf(to)

	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("reflect.Copy expects pointer in second argument, %s given", t.Kind())
	}

	if from == nil && isNilable(t.Elem().Kind()) {
		f = reflect.Zero(t.Elem().Type())
	}

	t.Elem().Set(f)
	return nil
}
