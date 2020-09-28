package container

type paramContainer interface {
	GetParam(string) (interface{}, error)
	MustGetParam(string) interface{}
	RegisterParam(string, ParamDefinition) error
	OverrideParam(string, ParamDefinition)
	HasParam(string) bool
	GetAllParamIDs() []string
}
