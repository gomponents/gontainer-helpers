package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type transaction struct {
	id int
}

type userRepo struct {
	transaction *transaction
}

type itemRepo struct {
	transaction *transaction
}

type purchaseService struct {
	userRepo *userRepo
	itemRepo *itemRepo
}

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
			Scope: ScopeShared,
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
			Scope: ScopeShared,
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
			Scope: ScopeShared,
		})
		container.Override("db", ServiceDefinition{
			Provider: func() (i interface{}, e error) {
				return struct{}{}, nil
			},
			Scope: ScopeShared,
		})
		container.Override("holding", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return struct{}{}, nil
			},
			Scope: ScopeShared,
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

	t.Run("Given scope", func(t *testing.T) {
		newContainer := func(s Scope) *Container {
			transactionID := 0
			c := NewContainer(nil)
			c.Override("transaction", ServiceDefinition{
				Provider: func() (interface{}, error) {
					transactionID++
					return &transaction{id: transactionID}, nil
				},
				Scope: s,
			})
			c.Override("userRepo", ServiceDefinition{
				Provider: func() (interface{}, error) {
					return &userRepo{
						transaction: c.MustGet("transaction").(*transaction),
					}, nil
				},
				Scope: s,
			})
			c.Override("itemRepo", ServiceDefinition{
				Provider: func() (interface{}, error) {
					return &itemRepo{
						transaction: c.MustGet("transaction").(*transaction),
					}, nil
				},
				Scope: s,
			})
			c.Override("purchaseService", ServiceDefinition{
				Provider: func() (interface{}, error) {
					return &purchaseService{
						userRepo: c.MustGet("userRepo").(*userRepo),
						itemRepo: c.MustGet("itemRepo").(*itemRepo),
					}, nil
				},
				Scope: s,
			})
			return c
		}

		assertEqualValues := func(t *testing.T, first interface{}, vals ...interface{}) {
			for _, v := range vals {
				assert.Equal(t, first, v, append([]interface{}{first}, vals...))
			}
		}

		t.Run(ScopeShared.String(), func(t *testing.T) {
			c := newContainer(ScopeShared)
			ps1 := c.MustGet("purchaseService").(*purchaseService)
			ps2 := c.MustGet("purchaseService").(*purchaseService)
			assertEqualValues(
				t,
				1,
				ps1.userRepo.transaction.id,
				ps1.itemRepo.transaction.id,
				ps2.userRepo.transaction.id,
				ps2.itemRepo.transaction.id,
			)
		})
		t.Run(ScopeNestedShared.String(), func(t *testing.T) {
			c := newContainer(ScopeNestedShared)
			for i := 1; i <= 3; i++ {
				ps := c.MustGet("purchaseService").(*purchaseService)
				assertEqualValues(
					t,
					i,
					ps.userRepo.transaction.id,
					ps.itemRepo.transaction.id,
				)
			}
		})
		t.Run(ScopeNonShared.String(), func(t *testing.T) {
			c := newContainer(ScopeNonShared)
			for i := 0; i < 3; i++ {
				ps := c.MustGet("purchaseService").(*purchaseService)
				ctID := i*2 + 1
				assert.Equal(
					t,
					[]int{ctID, ctID + 1},
					[]int{ps.userRepo.transaction.id, ps.itemRepo.transaction.id},
				)
			}
		})
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
