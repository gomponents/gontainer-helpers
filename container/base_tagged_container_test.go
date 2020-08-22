package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseTaggedContainer_GetByTag(t *testing.T) {
	container := NewBaseContainer(map[string]ServiceDefinition{
		"foo": {
			Provider: func() (i interface{}, e error) {
				return "foo", nil
			},
			Disposable: false,
		},
		"bar": {
			Provider: func() (i interface{}, e error) {
				return "bar", nil
			},
			Disposable: false,
		},
		"foobar": {
			Provider: func() (i interface{}, e error) {
				return "foobar", nil
			},
			Disposable: false,
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
