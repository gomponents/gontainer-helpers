package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamContainer_GetParam(t *testing.T) {
	t.Run("Given errors", func(t *testing.T) {
		circularContainer := NewParamContainer(nil)
		deps := map[string]string{
			"username": "nickname",
			"nickname": "name",
			"name":     "username",
		}
		for k, dep := range deps {
			d := dep
			circularContainer.OverrideParam(k, ParamDefinition{
				Provider: func() (interface{}, error) {
					return circularContainer.GetParam(d)
				},
			})
		}

		c2 := NewParamContainer(nil)
		c2.OverrideParam("name", ParamDefinition{
			Provider: func() (interface{}, error) {
				return circularContainer.GetParam("name")
			},
		})

		scenarios := []struct {
			container *ParamContainer
			id        string
			error     string
		}{
			{
				container: NewParamContainer(map[string]ParamDefinition{
					"db.host": {
						Provider: func() (interface{}, error) {
							return nil, fmt.Errorf("todo")
						},
						Disposable: false,
					},
				}),
				id:    "db.host",
				error: "cannot get parameter `db.host`: todo",
			},
			{
				container: NewParamContainer(nil),
				id:        "db.host",
				error:     "parameter `db.host` does not exist",
			},
			{
				container: circularContainer,
				id:        "nickname",
				error:     "cannot get parameter `nickname`: circular dependency: nickname -> name -> username -> nickname",
			},
			{
				container: c2,
				id:        "name",
				error:     "cannot get parameter `name`: cannot get parameter `name`: circular dependency: name -> username -> nickname -> name",
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				p, err := s.container.GetParam(s.id)
				assert.Nil(t, p)
				assert.EqualError(t, err, s.error)
			})
		}
	})
}
