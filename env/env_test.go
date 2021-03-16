package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	scenarios := []struct {
		envs     map[string]string
		key      string
		def      []string
		expected string
		error    string
	}{
		{
			key:   "MYSQL_USER",
			error: "environment variable `MYSQL_USER` does not exist",
		},
		{
			key:      "MYSQL_USER",
			def:      []string{"root"},
			expected: "root",
		},
		{
			envs:     map[string]string{"MYSQL_USER": "mysql-user"},
			key:      "MYSQL_USER",
			expected: "mysql-user",
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

			v, err := Get(s.key, s.def...)

			if s.error == "" {
				assert.NoError(t, err)
				assert.Equal(
					t,
					s.expected,
					v,
				)
				return
			}

			assert.EqualError(t, err, s.error)
			assert.Equal(t, "", v)
		})
	}
}

func TestGetInt(t *testing.T) {
	t.Run("Given scenarios", func(t *testing.T) {
		scenarios := []struct {
			envs     map[string]string
			key      string
			def      []int
			expected int
			error    string
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

				if s.error == "" {
					assert.NoError(t, err)
					assert.Equal(
						t,
						s.expected,
						v,
					)
					return
				}

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
