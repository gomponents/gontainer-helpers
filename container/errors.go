package container

import (
	"fmt"
	"strings"
)

type finalErr struct {
	error
}

func newCircularDepError(deps []string) finalErr {
	return finalErr{error: fmt.Errorf(
		"circular dependency: %s",
		strings.Join(deps, " -> "),
	)}
}
