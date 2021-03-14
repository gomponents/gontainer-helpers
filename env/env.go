package env

import (
	"fmt"
	"os"
	"strconv"
)

const envVarDoesntExist = "environment variable `%s` does not exist"

// Get returns environment variable if exists,
// otherwise it returns second argument if given.
func Get(key string, def ...string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0], nil
		}
		return "", fmt.Errorf(envVarDoesntExist, key)
	}
	return val, nil
}

// GetInt returns environment variable converted to int if exists,
// otherwise it returns second argument if given.
func GetInt(key string, def ...int) (int, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0], nil
		}
		return 0, fmt.Errorf(envVarDoesntExist, key)
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("cannot cast env(`%s`) to int: %s", key, err.Error())
	}
	return res, nil
}
