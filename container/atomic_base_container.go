package container

import (
	"sync"
)

type AtomicBaseContainer struct {
	*BaseContainer
	mutex *sync.Mutex
}

func NewAtomicBaseContainer(container *BaseContainer) *AtomicBaseContainer {
	return &AtomicBaseContainer{
		BaseContainer: container,
		mutex:         &sync.Mutex{},
	}
}

func (a AtomicBaseContainer) Get(id string) (interface{}, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.BaseContainer.Get(id)
}

func (a AtomicBaseContainer) Register(id string, s ServiceDefinition) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Register(id, s)
}

func (a AtomicBaseContainer) Override(id string, s ServiceDefinition) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.BaseContainer.Override(id, s)
}

func (a AtomicBaseContainer) Has(id string) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.BaseContainer.Has(id)
}

func (a AtomicBaseContainer) GetAllServiceIDs() []string {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.BaseContainer.GetAllServiceIDs()
}
