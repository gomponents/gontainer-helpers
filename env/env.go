package env

import (
	"fmt"
	"os"
	"strconv"
)

// MustGetInt returns environment variable converted to int if exists,
// otherwise return second argument if given
func MustGetInt(key string, def ...int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		panic(fmt.Sprintf("environment variable `%s` does not exist", key))
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("cannot cast env(`%s`) to int: %s", key, err.Error()))
	}
	return res
}
