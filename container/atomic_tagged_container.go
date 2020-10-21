package container

import (
	"sync"
)

type taggedContainer interface {
	GetByTag(string) ([]interface{}, error)
	MustGetByTag(string) []interface{}
	TagService(id string, tag string, priority int) error
	OverrideTagService(id string, tag string, priority int)
	IsTaggedBy(id string, tag string) bool
}

type AtomicTaggedContainer struct {
	container taggedContainer
	locker    sync.Locker
}

func NewAtomicTaggedContainer(container taggedContainer) *AtomicTaggedContainer {
	return &AtomicTaggedContainer{
		container: container,
		locker:    &sync.Mutex{},
	}
}

func (a AtomicTaggedContainer) GetByTag(id string) ([]interface{}, error) {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.GetByTag(id)
}

func (a AtomicTaggedContainer) MustGetByTag(id string) []interface{} {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.MustGetByTag(id)
}

func (a AtomicTaggedContainer) TagService(id string, tag string, priority int) error {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.TagService(id, tag, priority)
}

func (a AtomicTaggedContainer) OverrideTagService(id string, tag string, priority int) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.OverrideTagService(id, tag, priority)
}

func (a AtomicTaggedContainer) IsTaggedBy(id string, tag string) bool {
	a.locker.Lock()
	defer a.locker.Unlock()
	return a.container.IsTaggedBy(id, tag)
}
