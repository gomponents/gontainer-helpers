package container

import (
	"fmt"
	"sort"
)

type ServiceDefinition struct {
	Provider Provider

	// Disposable says whether object should be cached or no.
	Disposable bool

	// todo
	// Consistent says whether all dependencies should be shared even if they are disposable
	// see Container.GetConsistent
	Consistent bool
}

type metaServiceDefinition struct {
	definition ServiceDefinition
	service    interface{}
	created    bool
}

type Container struct {
	services           map[string]metaServiceDefinition
	circularDeps       *circularDeps
	decorators         []Decorator
	cacheGetConsistent map[string]interface{}
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

	if c.cacheGetConsistent != nil {
		if s, ok := c.cacheGetConsistent[id]; ok {
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

	if !serviceDef.definition.Disposable {
		serviceDef.created = true
		serviceDef.service = service
		c.services[id] = serviceDef
	}

	if c.cacheGetConsistent != nil {
		c.cacheGetConsistent[id] = service
	}

	return service, nil
}

// Revoke remove a cached copy of the service
func (c *Container) Revoke(id string) error {
	if !c.Has(id) {
		return fmt.Errorf("cannot revoke service `%s`, because it does not exist", id)
	}

	if c.services[id].definition.Disposable {
		return fmt.Errorf("cannot revoke service `%s`, because it is disposable", id)
	}

	if !c.services[id].created {
		return fmt.Errorf("cannot revoke service `%s`, because it is not created yet", id)
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
		return fmt.Errorf("cannot remove service `%s`, becaus it does not exist", id)
	}

	delete(c.services, id)

	return nil
}

func (c *Container) MustRemove(id string) {
	if err := c.Remove(id); err != nil {
		panic(err)
	}
}

// GetConsistent returns list of services at once.
// If two services use the same dependency, and the given dependency is disposable,
// both of them receive the same instance of given disposable dependency.
// In the following example userRepo and itemRepo will share the same transaction.
//
// c := NewContainer(nil)
// c.Override("transaction", ServiceDefinition{
//     Provider: func() (interface{}, error) {
//         return c.MustGet("db").(*sql.DB).Begin()
//     },
//     Disposable: true,
// })
// c.Override("userRepo", ServiceDefinition{
//     Provider: func() (interface{}, error) {
//         return NewUserRepo(c.MustGet("transaction").(*sql.Tx)), nil
//     },
//     Disposable: true,
// })
// c.Override("itemRepo", ServiceDefinition{
//     Provider: func() (interface{}, error) {
//         return NewItemRepo(c.MustGet("transaction").(*sql.Tx)), nil
//     },
//     Disposable: true,
// })
// services, err := c.GetConsistent("userRepo", "itemRepo", "transaction")
// // .. some logic
// services["transaction"].(*sql.Tx).Commit()
func (c *Container) GetConsistent(ids ...string) (map[string]interface{}, error) {
	c.cacheGetConsistent = make(map[string]interface{})
	defer func() {
		c.cacheGetConsistent = nil
	}()

	r := make(map[string]interface{})

	for _, id := range ids {
		var err error
		r[id], err = c.Get(id)
		if err != nil {
			return nil, fmt.Errorf("GetConsistent: %s", err.Error())
		}
	}

	return r, nil
}

func (c *Container) MustGetConsistent(ids ...string) map[string]interface{} {
	r, err := c.GetConsistent(ids...)
	if err != nil {
		panic(err)
	}
	return r
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
