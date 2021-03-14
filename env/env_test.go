package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInt(t *testing.T) {
	t.Run("Given scenarios", func(t *testing.T) {
		scenarios := []struct {
			envs     map[string]string
			key      string
			def      []int
			expected int
		}{
			{
				envs:     map[string]string{"GONTAINER_KEY": "123"},
				key:      "GONTAINER_KEY",
				expected: 123,
			},
			{
				envs:     map[string]string{"GONTAINER_KEY": "123"},
				key:      "GONTAINER_KEY2",
				def:      []int{500},
				expected: 500,
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				// prevent failure in case of existing environment variables
				unsetEnvs(t, s.envs)
				unsetEnvs(t, map[string]string{s.key: ""})

				setEnvs(t, s.envs)
				defer func() {
					unsetEnvs(t, s.envs)
				}()

				v, err := GetInt(s.key, s.def...)
				assert.NoError(t, err)
				assert.Equal(
					t,
					s.expected,
					v,
				)
			})
		}
	})

	t.Run("Given errors", func(t *testing.T) {
		scenarios := []struct {
			envs  map[string]string
			key   string
			def   []int
			error string
		}{
			{
				envs:  map[string]string{"GONTAINER_KEY": "hello"},
				key:   "GONTAINER_KEY",
				error: "cannot cast env(`GONTAINER_KEY`) to int: strconv.Atoi: parsing \"hello\": invalid syntax",
			},
			{
				envs:  map[string]string{"GONTAINER_KEY": "hello"},
				key:   "GONTAINER_KEY2",
				error: "environment variable `GONTAINER_KEY2` does not exist",
			},
			{
				envs:  map[string]string{"GONTAINER_KEY": "1.2"},
				key:   "GONTAINER_KEY",
				error: "cannot cast env(`GONTAINER_KEY`) to int: strconv.Atoi: parsing \"1.2\": invalid syntax",
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				// prevent failure in case of existing environment variables
				unsetEnvs(t, s.envs)
				unsetEnvs(t, map[string]string{s.key: ""})

				setEnvs(t, s.envs)
				defer func() {
					unsetEnvs(t, s.envs)
				}()

				v, err := GetInt(s.key, s.def...)
				assert.EqualError(t, err, s.error)
				assert.Equal(t, 0, v)
			})
		}
	})
}

func unsetEnvs(t *testing.T, vars map[string]string) {
	for k, _ := range vars {
		assert.NoError(
			t,
			os.Unsetenv(k),
		)
	}
}

func setEnvs(t *testing.T, vars map[string]string) {
	for k, v := range vars {
		assert.NoError(
			t,
			os.Setenv(k, v),
		)
	}
}
