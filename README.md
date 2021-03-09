[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer-helpers/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer-helpers?branch=master)

TEST

### API

#### Setters

Setters package allows easily to append a value to the field of struct by name. It supports unexported fields as well.

```go
package main

import (
	"fmt"

	"github.com/gomponents/gontainer-helpers/setter"
)

type Person struct {
	Name string
}

func main() {
	p := Person{}
	setter.Set(&p, "Name", "Jane")
	fmt.Println(p.Name) // Jane
}
```

#### Exporters

Exporters package allows exporting variables to executable GO code.

```go
package main

import (
	"fmt"

	"github.com/gomponents/gontainer-helpers/exporters"
)

func main() {
	v, _ := exporters.Export([]int{1, 2, 3})
	fmt.Println(v) // []int{int(1), int(2), int(3)}

	// ToString casts input value to string. It supports bool, nil and numeric values.
	s, _ := exporters.ToString(3.14) // s == "3.14"

	// panic: "cannot cast parameter of type `struct {}` to string: parameter of type `struct {}` is not supported"
	s2 := exporters.MustToString(struct{}{})
}
```
