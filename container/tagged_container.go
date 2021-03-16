package container

import (
	"fmt"
	"sort"
)

type TaggedContainer struct {
	container interface {
		Get(string) (interface{}, error)
	}
	mapping map[string]map[string]int // mapping[tag][serviceID] = priority
}

func NewTaggedContainer(
	container interface {
		Get(string) (interface{}, error)
	},
) *TaggedContainer {
	return &TaggedContainer{
		container: container,
		mapping:   make(map[string]map[string]int),
	}
}

// todo export?
func (b TaggedContainer) getIDsByTag(tag string) []string {
	svcs := make([]struct {
		id       string
		priority int
	}, 0)

	tagMapping, _ := b.mapping[tag]

	for id, p := range tagMapping {
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

	r := make([]string, len(tagMapping))
	for i, s := range svcs {
		r[i] = s.id
	}
	return r
}

func (b TaggedContainer) GetByTag(tag string) ([]interface{}, error) {
	ids := b.getIDsByTag(tag)
	result := make([]interface{}, len(ids))
	for i, id := range b.getIDsByTag(tag) {
		r, err := b.container.Get(id)
		if err != nil {
			return nil, fmt.Errorf("cannot get services by tag `%s`: %s", tag, err.Error())
		}
		result[i] = r
	}

	return result, nil
}

func (b TaggedContainer) MustGetByTag(tag string) []interface{} {
	result, err := b.GetByTag(tag)
	if err != nil {
		panic(err)
	}
	return result
}

func (b TaggedContainer) TagService(id string, tag string, priority int) error {
	if _, ok := b.mapping[tag]; !ok {
		b.mapping[tag] = make(map[string]int)
	}

	if _, ok := b.mapping[tag][id]; ok {
		return fmt.Errorf("service `%s` is already tagged as `%s`", id, tag)
	}

	b.mapping[tag][id] = priority

	return nil
}

func (b TaggedContainer) OverrideTagService(id string, tag string, priority int) {
	if _, ok := b.mapping[tag]; !ok {
		b.mapping[tag] = make(map[string]int)
	}

	b.mapping[tag][id] = priority
}

func (b TaggedContainer) IsTaggedBy(id string, tag string) bool {
	_, ok := b.mapping[tag][id]
	return ok
}
