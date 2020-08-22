package container

import (
	"fmt"
	"sort"
)

type BaseTaggedContainer struct {
	container Container
	mapping   map[string]map[string]int
}

func NewBaseTaggedContainer(container Container) *BaseTaggedContainer {
	return &BaseTaggedContainer{
		container: container,
		mapping:   make(map[string]map[string]int), // mapping[tag][serviceID] = priority
	}
}

func (b BaseTaggedContainer) GetByTag(tag string) ([]interface{}, error) {
	svcs := make([]struct {
		id       string
		priority int
	}, 0)

	for id, p := range b.mapping[tag] {
		svcs = append(
			svcs,
			struct {
				id       string
				priority int
			}{id: id, priority: p},
		)
	}

	sort.SliceStable(svcs, func(i, j int) bool {
		return svcs[i].priority > svcs[j].priority
	})

	result := make([]interface{}, 0)
	for _, s := range svcs {
		r, err := b.container.Get(s.id)
		if err != nil {
			return nil, fmt.Errorf("cannot get services by tag `%s`: %s", tag, err.Error())
		}
		result = append(result, r)
	}

	return result, nil
}

func (b BaseTaggedContainer) MustGetByTag(tag string) []interface{} {
	result, err := b.GetByTag(tag)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (b BaseTaggedContainer) TagService(tag string, id string, priority int) error {
	if _, ok := b.mapping[tag]; !ok {
		b.mapping[tag] = make(map[string]int)
	}

	if _, ok := b.mapping[tag][id]; ok {
		return fmt.Errorf("service `%s` is already tagged as `%s`", id, tag)
	}

	b.mapping[tag][id] = priority

	return nil
}

func (b BaseTaggedContainer) OverrideTagService(tag string, id string, priority int) {
	if _, ok := b.mapping[tag]; !ok {
		b.mapping[tag] = make(map[string]int)
	}

	b.mapping[tag][id] = priority
}
