package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAtomicParamContainer(t *testing.T) {
	t.Run("NewAtomicParamContainer", func(t *testing.T) {
		base := NewParamContainer(nil)
		c := NewAtomicParamContainer(base)

		assert.Same(t, base, c.container)
		assert.NotNil(t, c.locker)
	})

	t.Run("Concurrency", func(t *testing.T) {
		c := NewAtomicParamContainer(NewParamContainer(nil))
		g := goroutineGroup{}

		for i := 0; i < 100; i++ {
			id := fmt.Sprintf("param-%d", i)
			g.Go(func() {
				err := c.RegisterParam(id, ParamDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
				if err != nil {
					assert.EqualError(
						t,
						err,
						fmt.Sprintf("parameter `%s` already exists", id),
					)
				}
			})
			g.Go(func() {
				c.OverrideParam(id, ParamDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
			g.Go(func() {
				_, err := c.GetParam(id)
				if err != nil {
					assert.EqualError(t, err, fmt.Sprintf("parameter `%s` does not exist", id))
				}
			})
			g.Go(func() {
				defer func() {
					r := recover()
					if r != nil {
						assert.Equal(
							t,
							fmt.Sprintf("parameter `%s` does not exist", id),
							fmt.Sprintf("%s", r),
						)
					}
				}()
				c.MustGetParam(id)
			})
			g.Go(func() {
				c.HasParam(id)
			})
			g.Go(func() {
				c.GetAllParamIDs()
			})
		}

		g.Wait()
	})
}
