package container

import (
	"reflect"
	"testing"
)

func listMethods(o interface{}) []string {
	r := reflect.TypeOf(o)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	m := make([]string, r.NumMethod())
	for i := 0; i < r.NumMethod(); i++ {
		m[i] = r.Method(i).Name
	}
	return m
}

func assertNoAmbiguousMethods(t *testing.T, a, b interface{}) {
	l := make(map[string]bool)

	for _, m := range listMethods(a) {
		l[m] = true
	}

	for _, m := range listMethods(b) {
		if _, exists := l[m]; !exists {
			continue
		}
		t.Errorf("cannot composite `%T` and `%T`: ambigious method `%s`", a, b, m)
	}
}

func assertNoAmbiguousMethodsVariadic(t *testing.T, o ...interface{}) {
	for i, o1 := range o {
		for j, o2 := range o {
			if i == j || i > j {
				continue
			}
			assertNoAmbiguousMethods(t, o1, o2)
		}
	}
}

func TestComposition(t *testing.T) {
	// I want to give an option to do:
	//
	// type MyContainer struct {
	//     Container
	//     ParamContainer
	// }
	// c := MyContainer{}
	// c.Get("db")
	//
	// The following test checks whether do we have any ambiguous method names.
	t.Run("No ambiguous methods in interfaces", func(t *testing.T) {
		assertNoAmbiguousMethodsVariadic(
			t,
			(*container)(nil),
			(*paramContainer)(nil),
			(*taggedContainer)(nil),
		)
	})
	t.Run("No ambiguous methods in structs", func(t *testing.T) {
		assertNoAmbiguousMethodsVariadic(
			t,
			Container{},
			ParamContainer{},
			TaggedContainer{},
		)
	})
	t.Run("No ambiguous methods in atomic structs", func(t *testing.T) {
		assertNoAmbiguousMethodsVariadic(
			t,
			AtomicContainer{},
			AtomicParamContainer{},
			AtomicTaggedContainer{},
		)
	})
}
