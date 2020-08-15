package std

// GetMissingParameter panics always. It's just a mock param provider.
func GetMissingParameter(params ...string) interface{} {
	if len(params) > 0 {
		panic(params[0])
	}
	panic("missing parameter")
}
