package container

import (
	"sync"
)

type paramContainer interface {
	GetParam(string) (interface{}, error)
	MustGetParam(string) interface{}
	RegisterParam(string, ParamDefinition) error
	OverrideParam(string, ParamDefinition)
	HasParam(string) bool
	GetAllParamIDs() []string
}

type AtomicParamContainer struct {
	container paramContainer
	locker    sync.Locker
}

func NewAtomicParamContainer(container paramContainer) *AtomicParamContainer {
	return &AtomicParamContainer{
		container: container,
		locker:    &sync.Mutex{},
	}
}

func (a AtomicParamContainer) GetParam(id string) (interface{}, error) {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.GetParam(id)
}

func (a AtomicParamContainer) MustGetParam(id string) interface{} {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.MustGetParam(id)
}

func (a AtomicParamContainer) RegisterParam(id string, d ParamDefinition) error {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.RegisterParam(id, d)
}

func (a AtomicParamContainer) OverrideParam(id string, d ParamDefinition) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.container.OverrideParam(id, d)
}

func (a AtomicParamContainer) HasParam(id string) bool {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.HasParam(id)
}

func (a AtomicParamContainer) GetAllParamIDs() []string {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.GetAllParamIDs()
}
