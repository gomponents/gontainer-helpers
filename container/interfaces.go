package container

//type Container interface {
//	Get(string) (interface{}, error)
//	MustGet(string) interface{}
//	Has(string) bool
//}

type TaggedContainer interface {
	GetByTag(string) ([]interface{}, error)
	MustGetByTag(string) []interface{}
}

type ParamContainer interface {
	GetParam(string) (interface{}, error)
	MustGetParam(string) interface{}
	HasParam(string) bool
}
