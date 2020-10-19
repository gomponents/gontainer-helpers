package container

import (
	"sync"
)

type mockContainer struct {
	has     bool
	service interface{}
	error   error
}

func (m mockContainer) Get(string) (interface{}, error) {
	return m.service, m.error
}

type goroutineGroup struct {
	wg sync.WaitGroup
}

func (g *goroutineGroup) Go(f func()) {
	g.wg.Add(1)
	go func() {
		f()
		g.wg.Done()
	}()
}

func (g *goroutineGroup) Wait() {
	g.wg.Wait()
}
