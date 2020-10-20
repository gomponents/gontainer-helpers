package container

import (
	"fmt"
	"sort"
)

type ParamContainer struct {
	params       map[string]metaParamDefinition
	circularDeps *circularDeps
}

// todo replace by Provider
type ParamProvider func() interface{}

type ParamDefinition struct {
	Provider   ParamProvider
	Disposable bool
}

type metaParamDefinition struct {
	definition ParamDefinition
	param      interface{}
	created    bool
}

func NewParamContainer(definitions map[string]ParamDefinition) *ParamContainer {
	meta := make(map[string]metaParamDefinition)
	for n, v := range definitions {
		meta[n] = metaParamDefinition{
			definition: v,
			param:      nil,
			created:    false,
		}
	}

	return &ParamContainer{
		params:       meta,
		circularDeps: newCircularDeps(),
	}
}

func (b ParamContainer) GetParam(id string) (param interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cannot get parameter `%s`: %s", id, r)
		}
	}()

	defer b.circularDeps.stop()
	if deps := b.circularDeps.start(id); len(deps) != 0 {
		return nil, newCircularDepError(deps)
	}

	if !b.HasParam(id) {
		return nil, fmt.Errorf("parameter `%s` does not exist", id)
	}

	paramDef := b.params[id]
	if paramDef.created {
		return paramDef.param, nil
	}

	param = paramDef.definition.Provider()
	//if err != nil {
	//	if finalErr, ok := err.(finalErr); ok {
	//		return nil, finalErr
	//	}
	//	return nil, fmt.Errorf("cannot get parameter `%s`: %s", id, err.Error())
	//}

	if !paramDef.definition.Disposable {
		paramDef.created = true
		paramDef.param = param
	}

	return param, nil
}

func (b ParamContainer) RegisterParam(id string, d ParamDefinition) error {
	if b.HasParam(id) {
		return fmt.Errorf("parameter `%s` already exists", id)
	}
	b.params[id] = metaParamDefinition{
		definition: d,
		param:      nil,
		created:    false,
	}
	return nil
}

func (b ParamContainer) OverrideParam(id string, d ParamDefinition) {
	b.params[id] = metaParamDefinition{
		definition: d,
		param:      nil,
		created:    false,
	}
}

func (b ParamContainer) MustGetParam(id string) interface{} {
	r, err := b.GetParam(id)
	if err != nil {
		panic(err)
	}
	return r
}

func (b ParamContainer) HasParam(id string) bool {
	_, ok := b.params[id]
	return ok
}

func (b ParamContainer) GetAllParamIDs() []string {
	r := make([]string, 0)
	for n, _ := range b.params {
		r = append(r, n)
	}
	sort.Strings(r)
	return r
}
