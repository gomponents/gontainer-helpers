package container

import (
	"sync"
)

type AtomicContainer struct {
	*Container
	mutex *sync.Mutex
}

func NewAtomicContainer(container *Container) *AtomicContainer {
	return &AtomicContainer{
		Container: container,
		mutex:     &sync.Mutex{},
	}
}

func (a AtomicContainer) Get(id string) (interface{}, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Container.Get(id)
}

func (a AtomicContainer) Register(id string, s ServiceDefinition) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Register(id, s)
}

func (a AtomicContainer) Override(id string, s ServiceDefinition) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Container.Override(id, s)
}

func (a AtomicContainer) Has(id string) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Container.Has(id)
}

func (a AtomicContainer) GetAllServiceIDs() []string {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Container.GetAllServiceIDs()
}
