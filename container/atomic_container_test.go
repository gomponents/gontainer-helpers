package container

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAtomicContainer(t *testing.T) {
	base := NewContainer(nil)
	c := NewAtomicContainer(base)

	assert.Same(t, base, c.container)
	assert.NotNil(t, c.locker)
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

	// fatal error: concurrent map read and map write
	t.Run("MustGet and Override", func(t *testing.T) {
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
				defer func() {
					r := recover()
					if r != nil {
						assert.Equal(
							t,
							fmt.Sprintf("service `%s` does not exist", id),
							fmt.Sprintf("%s", r),
						)
					}
				}()
				c.MustGet(id)
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

	t.Run("RegisterDecorator", func(t *testing.T) {
		base := NewContainer(nil)
		c := NewAtomicContainer(base)
		g := goroutineGroup{}
		for i := 0; i < max; i++ {
			g.Go(func() {
				c.RegisterDecorator(func(_ string, svc interface{}) (interface{}, error) {
					return svc, nil
				})
			})
		}
		g.Wait()
		assert.Equal(t, max, len(base.decorators))
	})

	t.Run("All", func(t *testing.T) {
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

			g.Go(func() {
				defer func() {
					r := recover()
					if r != nil {
						assert.Equal(t, fmt.Sprintf("service `%s` does not exist", id), fmt.Sprintf("%s", r))
					}
				}()
				c.MustGet(id)
			})

			g.Go(func() {
				c.RegisterDecorator(func(_ string, svc interface{}) (interface{}, error) {
					return svc, nil
				})
			})
		}
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

		l := make(chan bool, 1)
		defer close(l)
		l <- true

		go func() {
			v, _ = c.Get("alias")
			<-l
		}()

		select {
		case l <- true:
			assert.Fail(t, "Lock should not be released")
			<-l
		case <-time.After(delay):
		}

		assert.Nil(t, v)
	})

	t.Run("Nested dependency must be registered in subcontainer", func(t *testing.T) {
		sub := NewContainer(nil)
		c := NewAtomicContainer(sub)
		registerNested(sub)
		var v interface{}

		l := make(chan bool, 1)
		defer close(l)
		l <- true

		go func() {
			v, _ = c.Get("alias")
			<-l
		}()

		select {
		case l <- true:
			<-l
		case <-time.After(delay):
			assert.Fail(t, "Timeout should not occur")
		}

		assert.NotNil(t, v)
	})
}
