package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaggedContainer_GetByTag(t *testing.T) {
	t.Run("Given scenario", func(t *testing.T) {
		container := NewContainer(map[string]ServiceDefinition{
			"foo": {
				Provider: func() (interface{}, error) {
					return "foo", nil
				},
			},
			"bar": {
				Provider: func() (interface{}, error) {
					return "bar", nil
				},
			},
			"foobar": {
				Provider: func() (interface{}, error) {
					return "foobar", nil
				},
			},
		})

		tagged := NewTaggedContainer(container)
		assert.NoError(t, tagged.TagService("foobar", "foos", 300))
		assert.NoError(t, tagged.TagService("foo", "foos", 30))
		assert.NoError(t, tagged.TagService("bar", "foos", 500))

		foos, errFoos := tagged.GetByTag("foos")
		assert.NoError(t, errFoos)
		assert.Equal(
			t,
			[]interface{}{"bar", "foobar", "foo"},
			foos,
		)

		loggers, errLogger := tagged.GetByTag("logger")
		assert.NoError(t, errLogger)
		assert.Equal(
			t,
			[]interface{}{},
			loggers,
		)
	})

	t.Run("Given errors", func(t *testing.T) {
		t.Run("Service already tagged", func(t *testing.T) {
			c := NewTaggedContainer(NewContainer(nil))
			tagSvc := func() error {
				return c.TagService("cmd", "commandHelp", 100)
			}
			assert.NoError(t, tagSvc())
			assert.EqualError(t, tagSvc(), "service `cmd` is already tagged as `commandHelp`")
		})
		t.Run("Parent container returns error", func(t *testing.T) {
			c := NewTaggedContainer(mockContainer{
				error: fmt.Errorf("service does not exist"),
			})
			c.OverrideTagService("mysql", "db", 0)
			_, err := c.GetByTag("db")
			assert.EqualError(t, err, "cannot get services by tag `db`: service does not exist")
		})
	})
}

func TestTaggedContainer_IsTaggedBy(t *testing.T) {
	base := NewContainer(map[string]ServiceDefinition{
		"userRepository": {
			Provider: func() (interface{}, error) {
				return struct{}{}, nil
			},
		},
		"productRepository": {
			Provider: func() (interface{}, error) {
				return struct{}{}, nil
			},
		},
	})
	tagged := NewTaggedContainer(base)
	tagged.OverrideTagService("userRepository", "repo", 0)
	tagged.OverrideTagService("productRepository", "repo", 0)

	assert.True(t, tagged.IsTaggedBy("userRepository", "repo"))
	assert.True(t, tagged.IsTaggedBy("productRepository", "repo"))
	assert.False(t, tagged.IsTaggedBy("userRepository", "db"))
}

func TestTaggedContainer_MustGetByTag(t *testing.T) {
	t.Run("Given success", func(t *testing.T) {
		db := struct{}{}
		c := NewTaggedContainer(mockContainer{
			service: db,
		})
		c.OverrideTagService("mysql", "db", 0)
		assert.Equal(t, []interface{}{db}, c.MustGetByTag("db"))
	})
	t.Run("Given error", func(t *testing.T) {
		defer func() {
			r := recover()
			if !assert.Implements(t, (*error)(nil), r) {
				return
			}
			assert.EqualError(
				t,
				r.(error),
				"cannot get services by tag `db`: service `mysql` does not exists",
			)
		}()

		c := NewTaggedContainer(mockContainer{
			error: fmt.Errorf("service `mysql` does not exists"),
		})
		c.OverrideTagService("mysql", "db", 0)
		c.MustGetByTag("db")
	})
}
