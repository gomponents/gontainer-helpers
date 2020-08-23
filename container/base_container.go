package container

import (
	"fmt"
	"sort"
)

type Provider = func() (interface{}, error)

type ServiceDefinition struct {
	Provider   Provider
	Disposable bool
}

type metaServiceDefinition struct {
	definition ServiceDefinition
	service    interface{}
	created    bool
}

type BaseContainer struct {
	services     map[string]metaServiceDefinition
	circularDeps *circularDeps
}

type finalContainerErr struct {
	error
}

func NewBaseContainer(definitions map[string]ServiceDefinition) *BaseContainer {
	meta := make(map[string]metaServiceDefinition)
	for n, v := range definitions {
		meta[n] = metaServiceDefinition{
			definition: v,
			service:    nil,
			created:    false,
		}
	}

	return &BaseContainer{
		services:     meta,
		circularDeps: newCircularDeps(),
	}
}

// Register registers new service, returns error in when service already exists
func (b BaseContainer) Register(id string, s ServiceDefinition) error {
	if b.Has(id) {
		return fmt.Errorf("service `%s` is already registered", id)
	}

	b.services[id] = metaServiceDefinition{
		definition: s,
		service:    nil,
		created:    false,
	}
	return nil
}

// Override overrides or registers service
func (b BaseContainer) Override(id string, s ServiceDefinition) {
	b.services[id] = metaServiceDefinition{
		definition: s,
		service:    nil,
		created:    false,
	}
}

func (b BaseContainer) Get(id string) (service interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cannot create service `%s`: %s", id, r)
		}
	}()

	defer b.circularDeps.stop()
	if deps := b.circularDeps.start(id); deps != nil {
		return nil, newCircularDepError(deps)
	}

	if !b.Has(id) {
		return nil, fmt.Errorf("service `%s` does not exist", id)
	}

	serviceDef := b.services[id]
	if serviceDef.created {
		return serviceDef.service, nil
	}

	service, err = serviceDef.definition.Provider()

	if err != nil {
		if finalErr, ok := err.(finalContainerErr); ok {
			return nil, finalErr
		}
		return nil, fmt.Errorf("cannot create service `%s`: %s", id, err.Error())
	}

	if !serviceDef.definition.Disposable {
		serviceDef.created = true
		serviceDef.service = service
		b.services[id] = serviceDef
	}

	return service, nil
}

func (b BaseContainer) MustGet(id string) interface{} {
	r, e := b.Get(id)

	if e != nil {
		panic(e)
	}

	return r
}

func (b BaseContainer) Has(id string) bool {
	_, ok := b.services[id]
	return ok
}

func (b BaseContainer) GetAllServiceIDs() []string {
	r := make([]string, 0)
	for n, _ := range b.services {
		r = append(r, n)
	}
	sort.Strings(r)
	return r
}
