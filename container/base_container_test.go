package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseContainer_Get(t *testing.T) {
	t.Run("Given circular dependency", func(t *testing.T) {
		container := NewBaseContainer(nil)
		container.Override("company", ServiceDefinition{
			Provider: func() (i interface{}, e error) {
				emp, err := container.Get("employer")
				if err != nil {
					return nil, err
				}

				return struct {
					employer interface{}
				}{
					employer: emp,
				}, nil
			},
			Disposable: false,
		})
		container.Override("employer", ServiceDefinition{
			Provider: func() (i interface{}, e error) {
				company, err := container.Get("company")
				if err != nil {
					return nil, err
				}

				return struct {
					company interface{}
				}{
					company: company,
				}, nil
			},
			Disposable: false,
		})
		container.Override("management", ServiceDefinition{
			Provider: func() (i interface{}, e error) {
				company, err := container.Get("company")
				if err != nil {
					return nil, err
				}

				return struct {
					company interface{}
				}{
					company: company,
				}, nil
			},
			Disposable: false,
		})
		container.Override("db", ServiceDefinition{
			Provider: func() (i interface{}, e error) {
				return struct{}{}, nil
			},
			Disposable: false,
		})

		management, managementErr := container.Get("management")
		assert.Nil(t, management)
		assert.EqualError(
			t,
			managementErr,
			"circular dependency: management -> company -> employer -> company",
		)
		// make sure chain of dependencies is cleared
		assert.Empty(t, container.circularDeps.chain)

		_, dbErr := container.Get("db")
		assert.NoError(t, dbErr)
		assert.Empty(t, container.circularDeps.chain)
	})
}

func TestBaseContainer_RegisterDecorator(t *testing.T) {
	c := NewBaseContainer(nil)
	assert.Len(t, *c.decorators, 0)
	c.RegisterDecorator(func(_ string, i interface{}) (interface{}, error) {
		return i, nil
	})
	assert.Len(t, *c.decorators, 1)
}
