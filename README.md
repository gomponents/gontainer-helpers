[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer-helpers/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer-helpers?branch=master)

TEST

# Gontainer-helpers

Set of packages to easily build DI container based on external configuration.

## API

### Setters

Setters package allows easily to append a value to the field of struct by name. It supports unexported fields as well.

```go
package main

import (
	"fmt"

	"github.com/gomponents/gontainer-helpers/setter"
)

type Person struct {
	name string
}

func main() {
	p := Person{}
	setter.Set(&p, "name", "Jane")
	fmt.Println(p.name) // Jane
}
```

### Exporters

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

### Callers

Callers package allows calling given func with list of parameters without knowing types of them.

**Definitions**

1. *Provider* - func which returns one or two values, second value must be type of error if given.
2. *Wither* - method which creates copy of struct, overrides one field then return given copy.

**Examples**

```go
package main

import (
	"regexp"
)

type Person struct {
	Name string
}

type Matcher struct {
	regexp *regexp.Regexp
}

// provider
func NewPerson(name string) *Person {
	return &Person{Name: name}
}

// provider
func NewMatcher(expr string) (*Matcher, error) {
	var err error
	r := Matcher{}
	r.regexp, err = regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// wither
func (p Person) WithName(n string) Person {
	p.Name = n
	return p
}
```