[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer-helpers/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer-helpers?branch=master)

TEST

### API

#### Setters

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

```go
package main

import (
	"fmt"

	"github.com/gomponents/gontainer-helpers/exporters"
)

func main() {
	v, err := exporters.Export([]int{1, 2, 3})
	if err != nil {
		panic(err)
	}
	fmt.Println(v) // []int{int(1), int(2), int(3)}
}
```
