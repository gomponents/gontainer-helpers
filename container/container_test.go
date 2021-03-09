package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer_Get(t *testing.T) {
	t.Run("Given circular dependency", func(t *testing.T) {
		container := NewContainer(nil)
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
		container.Override("holding", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return struct{}{}, nil
			},
			Disposable: false,
		})
		container.RegisterDecorator(func(s string, i interface{}) (interface{}, error) {
			if s == "holding" {
				_, err := container.Get("company")
				if err != nil {
					return nil, err
				}
			}
			return i, nil
		})

		management, managementErr := container.Get("management")
		assert.Nil(t, management)
		assert.EqualError(
			t,
			managementErr,
			"cannot create service `management`: circular dependency: management -> company -> employer -> company",
		)
		// make sure chain of dependencies is cleared
		assert.Empty(t, container.circularDeps.chain)

		holding, holdingErr := container.Get("holding")
		assert.Nil(t, holding)
		assert.EqualError(
			t,
			holdingErr,
			"cannot decorate service `holding`: circular dependency: holding -> company -> employer -> company",
		)
		assert.Empty(t, container.circularDeps.chain)

		_, dbErr := container.Get("db")
		assert.NoError(t, dbErr)
		assert.Empty(t, container.circularDeps.chain)
	})
}

func TestContainer_RegisterDecorator(t *testing.T) {
	c := NewContainer(nil)
	assert.Len(t, *c.decorators, 0)
	c.RegisterDecorator(func(_ string, i interface{}) (interface{}, error) {
		return i, nil
	})
	assert.Len(t, *c.decorators, 1)
}
