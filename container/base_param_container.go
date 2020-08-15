package container

import (
	"fmt"
)

type BaseParamContainer struct {
	providers map[string]ParamProvider
	params    map[string]interface{}
}

type ParamProvider func() interface{}

func NewBaseParamContainer(providers map[string]ParamProvider) *BaseParamContainer {
	if providers == nil {
		providers = make(map[string]ParamProvider)
	}
	return &BaseParamContainer{
		providers: providers,
		params:    make(map[string]interface{}),
	}
}

func (b BaseParamContainer) GetParam(id string) (val interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("param `%s`: %s", id, r)
		}
	}()

	if b.HasParam(id) {
		param, ok := b.params[id]
		if !ok {
			param = b.providers[id]()
			b.params[id] = param
		}
		return param, nil
	}

	return nil, fmt.Errorf("parameter `%s` does not exist", id)
}

func (b BaseParamContainer) SetParam(id string, provider ParamProvider) error {
	if b.HasParam(id) {
		return fmt.Errorf("parameter `%s` already exists", id)
	}
	b.providers[id] = provider
	return nil
}

func (b BaseParamContainer) OverrideParam(id string, provider ParamProvider) {
	b.providers[id] = provider
	delete(b.params, id)
}

func (b BaseParamContainer) MustGetParam(id string) interface{} {
	r, err := b.GetParam(id)
	if err != nil {
		panic(err)
	}
	return r
}

func (b BaseParamContainer) HasParam(id string) bool {
	_, ok := b.providers[id]
	return ok
}
