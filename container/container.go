package container

import (
	"fmt"
	"sort"
)

type Provider func() (interface{}, error)
type Decorator func(string, interface{}) (interface{}, error)

type ServiceDefinition struct {
	Provider Provider
	// Disposable says whether object should be cached or no.
	Disposable bool
}

type metaServiceDefinition struct {
	definition ServiceDefinition
	service    interface{}
	created    bool
}

type Container struct {
	services     map[string]metaServiceDefinition
	circularDeps *circularDeps
	decorators   *[]Decorator
}

func NewContainer(definitions map[string]ServiceDefinition) *Container {
	meta := make(map[string]metaServiceDefinition)
	for n, v := range definitions {
		meta[n] = metaServiceDefinition{
			definition: v,
			service:    nil,
			created:    false,
		}
	}

	d := make([]Decorator, 0)

	return &Container{
		services:     meta,
		circularDeps: newCircularDeps(),
		decorators:   &d,
	}
}

// Register registers new service, returns error in when service already exists
func (b Container) Register(id string, s ServiceDefinition) error {
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
func (b Container) Override(id string, s ServiceDefinition) {
	//b.providersMutex.Lock()
	//defer b.providersMutex.Unlock()

	b.services[id] = metaServiceDefinition{
		definition: s,
		service:    nil,
		created:    false,
	}
}

func (b Container) Get(id string) (service interface{}, err error) {
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
		if finalErr, ok := err.(finalErr); ok {
			return nil, finalErr
		}
		return nil, fmt.Errorf("cannot create service `%s`: %s", id, err.Error())
	}

	service, err = b.decorate(id, service)
	if err != nil {
		if finalErr, ok := err.(finalErr); ok {
			return nil, finalErr
		}
		return nil, fmt.Errorf("cannot decorate service `%s`: %s", id, err.Error())
	}

	if !serviceDef.definition.Disposable {
		serviceDef.created = true
		serviceDef.service = service
		b.services[id] = serviceDef
	}

	return service, nil
}

func (b Container) MustGet(id string) interface{} {
	r, e := b.Get(id)

	if e != nil {
		panic(e)
	}

	return r
}

func (b Container) Has(id string) bool {
	_, ok := b.services[id]
	return ok
}

func (b Container) GetAllServiceIDs() []string {
	r := make([]string, 0)
	for n, _ := range b.services {
		r = append(r, n)
	}
	sort.Strings(r)
	return r
}

func (b Container) RegisterDecorator(d Decorator) {
	*b.decorators = append(*b.decorators, d)
}

func (b Container) decorate(id string, s interface{}) (r interface{}, err error) {
	r = s
	for _, d := range *b.decorators {
		r, err = d(id, r)
		if err != nil {
			return
		}
	}
	return
}
