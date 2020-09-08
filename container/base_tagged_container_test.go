package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseTaggedContainer_GetByTag(t *testing.T) {
	t.Run("Given scenario", func(t *testing.T) {
		container := NewBaseContainer(map[string]ServiceDefinition{
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

		tagged := NewBaseTaggedContainer(container)
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
			c := NewBaseTaggedContainer(NewBaseContainer(nil))
			tagSvc := func() error {
				return c.TagService("cmd", "commandHelp", 100)
			}
			assert.NoError(t, tagSvc())
			assert.EqualError(t, tagSvc(), "service `cmd` is already tagged as `commandHelp`")
		})
	})
}

func TestBaseTaggedContainer_IsTaggedBy(t *testing.T) {
	base := NewBaseContainer(map[string]ServiceDefinition{
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
	tagged := NewBaseTaggedContainer(base)
	tagged.OverrideTagService("userRepository", "repo", 0)
	tagged.OverrideTagService("productRepository", "repo", 0)

	assert.True(t, tagged.IsTaggedBy("userRepository", "repo"))
	assert.True(t, tagged.IsTaggedBy("productRepository", "repo"))
	assert.False(t, tagged.IsTaggedBy("userRepository", "db"))
}
