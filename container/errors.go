package container

import (
	"fmt"
	"strings"
)

type finalErr struct {
	error
}

func newCircularDepError(deps []string) finalErr {
	return finalErr{fmt.Errorf(
		"circular dependency: %s",
		strings.Join(deps, " -> "),
	)}
}
