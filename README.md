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
