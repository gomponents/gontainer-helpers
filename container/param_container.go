package container

import (
	"fmt"
	"sort"
)

type ParamContainer struct {
	params       map[string]metaParamDefinition
	circularDeps *circularDeps
}

// todo remove ParamDefinition, only Provider is enough
type ParamDefinition struct {
	Provider   Provider
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
	const errorMsg = "cannot get parameter `%s`: %s"
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if f, ok := r.(finalErr); ok {
			err = f
			if len(b.circularDeps.chain) == 0 {
				err = fmt.Errorf(errorMsg, id, err.Error())
			}
			return
		}
		err = fmt.Errorf(errorMsg, id, r)
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

	param, err = paramDef.definition.Provider()
	if err != nil {
		panic(err)
	}

	if !paramDef.definition.Disposable {
		paramDef.created = true
		paramDef.param = param
		b.params[id] = paramDef
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
