package container

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAtomicContainer(t *testing.T) {
	// todo
}

func TestAtomicContainer_Concurrency(t *testing.T) {
	const max = 100

	// fatal error: concurrent map writes
	t.Run("Override", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			g.Go(func() {
				c.Override("foo", ServiceDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
		}
		g.Wait()
	})

	// fatal error: concurrent map writes
	// fatal error: concurrent map read and map write
	t.Run("Register", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			id := fmt.Sprintf("svc-%d", i)
			g.Go(func() {
				_ = c.Register(id, ServiceDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
		}
		g.Wait()
	})

	// fatal error: concurrent map iteration and map write
	t.Run("GetAllServiceIDs and Override", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			id := fmt.Sprintf("svc-%d", i)
			g.Go(func() {
				c.GetAllServiceIDs()
			})
			g.Go(func() {
				c.Override(id, ServiceDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
		}
		g.Wait()
	})

	// fatal error: concurrent map read and map write
	t.Run("Get and Override", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			id := fmt.Sprintf("svc-%d", i)
			g.Go(func() {
				c.Override(id, ServiceDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
			g.Go(func() {
				_, _ = c.Get(id)
			})
		}
		g.Wait()
	})

	// fatal error: concurrent map iteration and map write
	t.Run("GetAllServiceIDs and Override", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			id := fmt.Sprintf("svc-%d", i)
			g.Go(func() {
				c.GetAllServiceIDs()
			})
			g.Go(func() {
				c.Override(id, ServiceDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
		}
	})

	// fatal error: concurrent map read and map write
	t.Run("Has and Override", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			id := fmt.Sprintf("svc-%d", i)
			g.Go(func() {
				c.Override(id, ServiceDefinition{
					Provider: func() (interface{}, error) {
						return nil, nil
					},
					Disposable: false,
				})
			})
			g.Go(func() {
				c.Has(id)
			})
		}
		g.Wait()
	})
}

func TestAtomicContainer_NestedLock(t *testing.T) {
	registerNested := func(c container) {
		c.Override("foo", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return struct{}{}, nil
			},
			Disposable: false,
		})
		c.Override("alias", ServiceDefinition{
			Provider: func() (interface{}, error) {
				return c.Get("foo")
			},
			Disposable: false,
		})
	}

	delay := 100 * time.Millisecond

	t.Run("Atomic does not support nested dependencies", func(t *testing.T) {
		c := NewAtomicContainer(NewContainer(nil))
		registerNested(c)
		var v interface{}
		go func() {
			v, _ = c.Get("alias")
		}()
		time.Sleep(delay)
		assert.Nil(t, v)
	})

	t.Run("Nested dependency must be registered in subcontainer", func(t *testing.T) {
		sub := NewContainer(nil)
		c := NewAtomicContainer(sub)
		registerNested(sub)
		var v interface{}
		go func() {
			v, _ = c.Get("alias")
		}()
		time.Sleep(delay)
		assert.NotNil(t, v)
	})
}
