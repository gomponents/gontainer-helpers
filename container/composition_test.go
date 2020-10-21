package container

import (
	"fmt"
	"reflect"
	"strings"
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
	t.Run("Print a composed interface", func(t *testing.T) {
		listMethods := func(o interface{}) []string {
			r := reflect.TypeOf(o)
			if r.Kind() == reflect.Ptr {
				r = r.Elem()
			}

			typeOf := func(t reflect.Type) string {
				r := t.Name()
				if t.PkgPath() != "" {
					r = fmt.Sprintf(`"%s".%s`, t.PkgPath(), r)
				}
				if r == "" {
					r = t.String()
				}
				return r
			}

			var methods []string

			for i := 0; i < r.NumMethod(); i++ {
				mr := r.Method(i)
				var inputs []string
				for j := 1; j < mr.Type.NumIn(); j++ {
					inputs = append(inputs, typeOf(mr.Type.In(j)))
				}

				var outputs []string
				for j := 0; j < mr.Type.NumOut(); j++ {
					outputs = append(outputs, typeOf(mr.Type.Out(j)))
				}
				out := ""
				switch len(outputs) {
				case 0:
				case 1:
					out = fmt.Sprintf(" %s", outputs[0])
				case 2:
					out = fmt.Sprintf(" (%s)", strings.Join(outputs, ", "))
				}

				method := fmt.Sprintf(
					"%s(%s)%s",
					mr.Name,
					strings.Join(inputs, ", "),
					out,
				)

				methods = append(methods, method)
			}

			return methods
		}

		m := []string{"// service container"}
		m = append(m, listMethods(AtomicContainer{})...)
		m = append(m, "\n", "// param container")
		m = append(m, listMethods(AtomicParamContainer{})...)
		m = append(m, "\n", "// tagged container")
		m = append(m, listMethods(AtomicTaggedContainer{})...)

		i := `type Container interface {
	%s
}`
		i = fmt.Sprintf(
			i,
			strings.Join(m, "\n\t"),
		)

		t.Log("\n" + i)
	})
}
