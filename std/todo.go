package std

const paramTodoMsg = "parameter todo"

// ParameterTodo panics always. It's just a mock param provider.
func ParameterTodo(params ...string) interface{} {
	if len(params) > 0 {
		panic(params[0])
	}
	panic(paramTodoMsg)
}
