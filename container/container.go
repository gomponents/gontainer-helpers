package container

import (
	"fmt"
	"sort"
)

type Scope uint

const (
	ScopeShared       Scope = iota // The same instance is used each time you request it from this container
	ScopeNestedShared              // The same instance is used only in the scope of the given service
	ScopeNonShared                 // New instance is created each time you request it from this container
)

func (s Scope) String() string {
	return []string{
		"ScopeShared",
		"ScopeNestedShared",
		"ScopeNonShared",
	}[s]
}

type ServiceDefinition struct {
	Provider Provider
	Scope    Scope
}

type metaServiceDefinition struct {
	definition ServiceDefinition
	service    interface{}
	created    bool
}

type Container struct {
	services          map[string]metaServiceDefinition
	circularDeps      *circularDeps
	decorators        []Decorator
	cacheNestedShared map[string]interface{}
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

	return &Container{
		services:     meta,
		circularDeps: newCircularDeps(),
		decorators:   make([]Decorator, 0),
	}
}

// Register registers new service, returns error in when service already exists
func (c *Container) Register(id string, s ServiceDefinition) error {
	if c.Has(id) {
		return fmt.Errorf("service `%s` is already registered", id)
	}

	c.services[id] = metaServiceDefinition{
		definition: s,
		service:    nil,
		created:    false,
	}
	return nil
}

// Override overrides or registers service
func (c *Container) Override(id string, s ServiceDefinition) {
	c.services[id] = metaServiceDefinition{
		definition: s,
		service:    nil,
		created:    false,
	}
}

func (c *Container) Get(id string) (service interface{}, err error) {
	const errorMsg = "cannot %s service `%s`: %s"

	defer func() {
		// todo check "if f, ok := r.(finalErr); ok {"
		// see ParamContainer.GetParam
		if r := recover(); r != nil {
			err = fmt.Errorf(errorMsg, "create", id, r)
		}
	}()

	defer c.circularDeps.stop()
	if deps := c.circularDeps.start(id); len(deps) != 0 {
		return nil, newCircularDepError(deps)
	}

	if c.cacheNestedShared != nil {
		if s, ok := c.cacheNestedShared[id]; ok {
			return s, nil
		}
	}

	if !c.Has(id) {
		return nil, fmt.Errorf("service `%s` does not exist", id)
	}

	serviceDef := c.services[id]
	if serviceDef.created {
		return serviceDef.service, nil
	}

	// c.cacheNestedShared == nil to avoid recreation empty map in consistent subdeps
	if c.cacheNestedShared == nil {
		c.cacheNestedShared = make(map[string]interface{})
		defer func() {
			c.cacheNestedShared = nil
		}()
	}

	decorateErr := func(err error, action string) error {
		if finalErr, ok := err.(finalErr); ok {
			if len(c.circularDeps.chain) == 1 {
				return fmt.Errorf(errorMsg, action, id, err.Error())
			}
			return finalErr
		}
		return fmt.Errorf(errorMsg, action, id, err.Error())
	}

	service, err = serviceDef.definition.Provider()
	if err != nil {
		return nil, decorateErr(err, "create")
	}

	service, err = c.decorate(id, service)
	if err != nil {
		return nil, decorateErr(err, "decorate")
	}

	if serviceDef.definition.Scope == ScopeShared {
		serviceDef.created = true
		serviceDef.service = service
		c.services[id] = serviceDef
	}

	if serviceDef.definition.Scope == ScopeNestedShared {
		c.cacheNestedShared[id] = service
	}

	return service, nil
}

const (
	errMsgServiceDoesNotExist = "cannot %s service `%s`, because it does not exist"
)

// Revoke remove a cached copy of the service
func (c *Container) Revoke(id string) error {
	if !c.Has(id) {
		return fmt.Errorf(errMsgServiceDoesNotExist, "revoke", id)
	}

	cp := c.services[id]
	cp.service = nil
	cp.created = false
	c.services[id] = cp

	return nil
}

func (c *Container) MustRevoke(id string) {
	if err := c.Revoke(id); err != nil {
		panic(err)
	}
}

// Remove removes service completely
func (c *Container) Remove(id string) error {
	if !c.Has(id) {
		return fmt.Errorf(errMsgServiceDoesNotExist, "remove", id)
	}

	delete(c.services, id)

	return nil
}

func (c *Container) MustRemove(id string) {
	if err := c.Remove(id); err != nil {
		panic(err)
	}
}

func (c *Container) MustGet(id string) interface{} {
	r, e := c.Get(id)

	if e != nil {
		panic(e)
	}

	return r
}

func (c *Container) Has(id string) bool {
	_, ok := c.services[id]
	return ok
}

func (c *Container) GetAllServiceIDs() []string {
	r := make([]string, 0)
	for n, _ := range c.services {
		r = append(r, n)
	}
	sort.Strings(r)
	return r
}

func (c *Container) RegisterDecorator(d Decorator) {
	c.decorators = append(c.decorators, d)
}

func (c *Container) decorate(id string, s interface{}) (r interface{}, err error) {
	r = s
	for _, d := range c.decorators {
		r, err = d(id, r)
		if err != nil {
			return
		}
	}
	return
}
