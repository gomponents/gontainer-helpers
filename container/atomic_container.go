package container

import (
	"sync"
)

type container interface {
	Get(string) (interface{}, error)
	MustGet(string) interface{}
	Register(string, ServiceDefinition) error
	Override(string, ServiceDefinition)
	Has(string) bool
	GetAllServiceIDs() []string
	RegisterDecorator(Decorator)
}

type AtomicContainer struct {
	container container
	mutex     *sync.Mutex
}

func NewAtomicContainer(c container) *AtomicContainer {
	return &AtomicContainer{
		container: c,
		mutex:     &sync.Mutex{},
	}
}

func (a AtomicContainer) Get(id string) (interface{}, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.container.Get(id)
}

func (a AtomicContainer) Register(id string, s ServiceDefinition) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.container.Register(id, s)
}

func (a AtomicContainer) Override(id string, s ServiceDefinition) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.container.Override(id, s)
}

func (a AtomicContainer) Has(id string) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.container.Has(id)
}

func (a AtomicContainer) GetAllServiceIDs() []string {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.container.GetAllServiceIDs()
}

func (a AtomicContainer) MustGet(id string) interface{} {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.container.MustGet(id)
}

func (a AtomicContainer) RegisterDecorator(d Decorator) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.container.RegisterDecorator(d)
}
