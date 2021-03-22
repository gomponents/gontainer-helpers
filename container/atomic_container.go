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
	Revoke(string) error
	MustRevoke(string)
	Remove(string) error
	MustRemove(string)
	GetSingletons(...string) (map[string]interface{}, error)
	MustGetSingletons(...string) map[string]interface{}
}

type AtomicContainer struct {
	container container
	locker    sync.Locker
}

func NewAtomicContainer(c container) *AtomicContainer {
	return &AtomicContainer{
		container: c,
		locker:    &sync.Mutex{},
	}
}

func (a AtomicContainer) Get(id string) (interface{}, error) {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.Get(id)
}

func (a AtomicContainer) Register(id string, s ServiceDefinition) error {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.Register(id, s)
}

func (a AtomicContainer) Override(id string, s ServiceDefinition) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.container.Override(id, s)
}

func (a AtomicContainer) Has(id string) bool {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.Has(id)
}

func (a AtomicContainer) GetAllServiceIDs() []string {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.GetAllServiceIDs()
}

func (a AtomicContainer) MustGet(id string) interface{} {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.MustGet(id)
}

func (a AtomicContainer) RegisterDecorator(d Decorator) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.container.RegisterDecorator(d)
}

func (a AtomicContainer) Revoke(id string) error {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.Revoke(id)
}

func (a AtomicContainer) MustRevoke(id string) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.container.MustRevoke(id)
}

func (a AtomicContainer) Remove(id string) error {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.Remove(id)
}

func (a AtomicContainer) MustRemove(id string) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.container.MustRemove(id)
}

func (a AtomicContainer) GetSingletons(ids ...string) (map[string]interface{}, error) {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.GetSingletons(ids...)
}

func (a AtomicContainer) MustGetSingletons(ids ...string) map[string]interface{} {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.MustGetSingletons(ids...)
}
