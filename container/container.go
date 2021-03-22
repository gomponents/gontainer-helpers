package container

import (
	"fmt"
	"sort"
)

type ServiceDefinition struct {
	Provider Provider

	// Disposable says whether object should be cached or no.
	// Disposable is logic negation of "shared" from Symfony.
	// https://symfony.com/doc/current/service_container/shared.html
	// todo remove
	Disposable bool

	// CollatedDeps says whether all sub-dependencies should be shared even if they are disposable.
	// see Container.GetConsistent
	// todo remove
	CollatedDeps bool

	// Singleton configures life-cycle of the given dependency
	// see https://docs.spring.io/spring-framework/docs/3.0.0.M3/reference/html/ch04s04.html
	// see https://symfony.com/doc/current/service_container/shared.html
	// todo use it
	Singleton bool

	// EnforceSingletonDeps changes life-cycle of sub-dependencies
	//
	// Let's consider the following example:
	// 1. PurchaseService depends on UserRepository and ItemRepository
	// 2. UserRepository depends on SQLTransaction
	// 3. ItemRepository depends on SQLTransaction
	// 4. PurchaseService, UserRepository, ItemRepository and SQLTransaction are not singletons
	//
	// Our dependency graph will look like:
	//
	// PurchaseService
	//               |-> UserRepository -> SQLTransaction (1)
	//               |-> ItemRepository -> SQLTransaction (2)
	//
	// In scope of one service we have 2 repositories. Both of them depends on different SQL transactions,
	// however we would like to achieve one transaction for entire scope of PurchaseService.
	// The given flag gives an option to re-use non-singletons in scope of one service.
	// When we enable it our dependency graph will look like:
	//
	// PurchaseService
	//               |-> UserRepository -> SQLTransaction
	//               |-> ItemRepository ↗
	//
	// Let's consider more complex scenario:
	// 1. PurchaseService is wrapped by PurchaseServiceSQLTransactionAware
	// 2. PurchaseServiceSQLTransactionAware depends on SQLTransaction
	//
	// PurchaseServiceSQLTransactionAware---------------------------------------|
	//                                  |-> PurchaseService                     ↓
	//                                                    |-> UserRepository -> SQLTransaction
	//                                                    |-> ItemRepository ↗
	//
	// PurchaseServiceSQLTransactionAware has the same method signature as PurchaseService.
	// However it does not perform any business logic, instead of that it calls method PurchaseService.DoAction
	// and depending on result it rollbacks or commits performed SQL operations.
	//
	// type PurchaseServiceSQLTransactionAware struct {
	//     purchaseService PurchaseService
	//     transaction     *sql.Tx
	// }
	//
	// func (p *PurchaseServiceSQLTransactionAware) DoAction() error {
	//     if err := p.purchaseService.DoAction(); err != nil {
	//         p.transaction.Rollback()
	//         return err
	//     }
	//
	//     p.transaction.Commit()
	//     return nil
	// }
	//
	// In the result your service is simpler. You do not need to handle your transaction in scope of business logic.
	// Instead of that you can wrap your business logic by transaction.
	// todo use it
	EnforceSingletonDeps bool
}

type metaServiceDefinition struct {
	definition ServiceDefinition
	service    interface{}
	created    bool
}

type Container struct {
	services        map[string]metaServiceDefinition
	circularDeps    *circularDeps
	decorators      []Decorator
	cacheGetCollate map[string]interface{}
	cacheGet        map[string]interface{}
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

	if c.cacheGetCollate != nil {
		if s, ok := c.cacheGetCollate[id]; ok {
			return s, nil
		}
	}

	if c.cacheGet != nil {
		if s, ok := c.cacheGet[id]; ok {
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

	// c.cacheGet == nil to avoid recreation empty map in consistent subdeps
	if c.cacheGet == nil && serviceDef.definition.CollatedDeps {
		c.cacheGet = make(map[string]interface{})
		defer func() {
			c.cacheGet = nil
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

	if !serviceDef.definition.Disposable {
		serviceDef.created = true
		serviceDef.service = service
		c.services[id] = serviceDef
	}

	if c.cacheGetCollate != nil {
		c.cacheGetCollate[id] = service
	}

	if c.cacheGet != nil {
		c.cacheGet[id] = service
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

// GetCollate returns list of services at once.
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
// services, err := c.GetCollate("userRepo", "itemRepo", "transaction")
// // .. some logic
// services["transaction"].(*sql.Tx).Commit()
func (c *Container) GetCollate(ids ...string) (map[string]interface{}, error) {
	c.cacheGetCollate = make(map[string]interface{})
	defer func() {
		c.cacheGetCollate = nil
	}()

	r := make(map[string]interface{})

	for _, id := range ids {
		var err error
		r[id], err = c.Get(id)
		if err != nil {
			return nil, fmt.Errorf("GetCollate: %s", err.Error())
		}
	}

	return r, nil
}

func (c *Container) MustGetCollate(ids ...string) map[string]interface{} {
	r, err := c.GetCollate(ids...)
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
