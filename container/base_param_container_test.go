package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseParamContainer_GetParam(t *testing.T) {
	t.Run("Given errors", func(t *testing.T) {
		scenarios := []struct {
			container *BaseParamContainer
			id        string
			error     string
		}{
			{
				container: NewBaseParamContainer(map[string]ParamProvider{
					"db.host": func() interface{} {
						panic("todo")
					},
				}),
				id:    "db.host",
				error: "cannot get parameter `db.host`: todo",
			},
			{
				container: NewBaseParamContainer(nil),
				id:        "db.host",
				error:     "parameter `db.host` does not exist",
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
