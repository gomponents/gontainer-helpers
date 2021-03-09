package panics

import (
	"fmt"
)

type Getter func() interface{}

func WrapGetter(g Getter, s string) interface{} {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		panic(fmt.Sprintf("%s: %s", s, r))
	}()
	return g()
}
