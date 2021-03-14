package std

import (
	"errors"
)

const paramTodoMsg = "parameter todo"

// ParameterTodo returns error always. It's just a mock param provider.
func ParameterTodo(params ...string) (interface{}, error) {
	if len(params) > 0 {
		return nil, errors.New(params[0])
	}
	return nil, errors.New(paramTodoMsg)
}
