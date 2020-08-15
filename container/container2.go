package container

import (
	"fmt"
)

type BaseTaggedContainer struct {
	container Container
	mapping   map[string][]string
}

func (b BaseTaggedContainer) Get(id string) (interface{}, error) {
	return b.container.Get(id)
}

func (b BaseTaggedContainer) MustGet(id string) interface{} {
	return b.container.MustGet(id)
}

func (b BaseTaggedContainer) Has(id string) bool {
	return b.container.Has(id)
}

func (b BaseTaggedContainer) GetByTag(tag string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, id := range b.mapping[tag] {
		s, e := b.container.Get(id)
		if e != nil {
			return nil, fmt.Errorf("cannot get services by tag %s due to: %s", tag, e.Error())
		}
		result = append(result, s)
	}
	return result, nil
}

func (b BaseTaggedContainer) MustGetByTag(tag string) []interface{} {
	r, e := b.GetByTag(tag)
	if e != nil {
		panic(e)
	}
	return r
}

func NewBaseTaggedContainer(container Container, mapping map[string][]string) *BaseTaggedContainer {
	return &BaseTaggedContainer{container: container, mapping: mapping}
}
