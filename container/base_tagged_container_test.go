package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseTaggedContainer_GetByTag(t *testing.T) {
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
	assert.NoError(t, tagged.TagService("foos", "foobar", 300))
	assert.NoError(t, tagged.TagService("foos", "foo", 30))
	assert.NoError(t, tagged.TagService("foos", "bar", 500))

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
	tagged.OverrideTagService("repo", "userRepository", 0)
	tagged.OverrideTagService("repo", "productRepository", 0)

	assert.True(t, tagged.IsTaggedBy("userRepository", "repo"))
	assert.True(t, tagged.IsTaggedBy("productRepository", "repo"))
	assert.False(t, tagged.IsTaggedBy("userRepository", "db"))
}
