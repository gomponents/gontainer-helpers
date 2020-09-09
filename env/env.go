package env

import (
	"fmt"
	"os"
	"strconv"
)

const envVarDoesntExist = "environment variable `%s` does not exist"

// MustGetInt returns environment variable if exists,
// otherwise it returns second argument if given.
func MustGet(key string, def ...string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		panic(fmt.Sprintf(envVarDoesntExist, key))
	}
	return val
}

// MustGetInt returns environment variable converted to int if exists,
// otherwise it returns second argument if given.
func MustGetInt(key string, def ...int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		panic(fmt.Sprintf(envVarDoesntExist, key))
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("cannot cast env(`%s`) to int: %s", key, err.Error()))
	}
	return res
}
