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
	assert.Len(t, c.decorators, 0)
	c.RegisterDecorator(func(_ string, i interface{}) (interface{}, error) {
		return i, nil
	})
	assert.Len(t, c.decorators, 1)
}

func TestContainer_GetSingletons(t *testing.T) {
	t.Run("Shared disposable dependency", func(t *testing.T) {
		c := NewContainer(nil)
		i := 0
		c.Override("inc", ServiceDefinition{
			Provider: func() (interface{}, error) {
				i++
				return i, nil
			},
			Disposable: true,
		})
		c.Override("sliceA", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return []interface{}{c.MustGet("inc")}, nil
			},
			Disposable: true,
		})
		c.Override("sliceB", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return []interface{}{c.MustGet("inc")}, nil
			},
			Disposable: true,
		})
		c.Override("metaslice", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return append(
					c.MustGet("sliceA").([]interface{}),
					c.MustGet("sliceB").([]interface{})...,
				), nil
			},
			Disposable:   true,
			CollatedDeps: true,
		})
		c.Override("metaslice2", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return append(
					c.MustGet("sliceA").([]interface{}),
					c.MustGet("sliceB").([]interface{})...,
				), nil
			},
			Disposable:   true,
			CollatedDeps: false,
		})

		slices, err := c.GetSingletons("sliceA", "sliceB")
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{1}, slices["sliceA"])
		assert.Equal(t, []interface{}{1}, slices["sliceB"])

		assert.Equal(t, []interface{}{2}, c.MustGet("sliceA"))
		assert.Equal(t, []interface{}{3}, c.MustGet("sliceB"))
		assert.Equal(t, []interface{}{4}, c.MustGet("sliceA"))

		slices, err = c.GetSingletons("sliceA", "sliceB")
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{5}, slices["sliceA"])
		assert.Equal(t, []interface{}{5}, slices["sliceB"])

		// todo better tests for CollatedDeps
		m := c.MustGet("metaslice").([]interface{})
		assert.Equal(t, m[0], m[1])
		m2 := c.MustGet("metaslice2").([]interface{})
		assert.NotEqual(t, m2[0], m2[1])

		assert.Nil(t, c.cacheGetSingletons)
		assert.Nil(t, c.cacheGet)
	})

	t.Run("Given error", func(t *testing.T) {
		c := NewContainer(nil)
		s, err := c.GetSingletons("db")
		assert.EqualError(t, err, "GetSingletons: service `db` does not exist")
		assert.Nil(t, s)
	})
}
