package caller

import (
	"fmt"
	"reflect"
)

type reflectType struct {
	reflect.Type
}

// inVariadicAware works almost same as reflect.Type.In,
// but it returns t.In(i).Elem() for t.isVariadic() && i >= t.NumIn().
func (t reflectType) inVariadicAware(i int) reflect.Type {
	last := t.NumIn() - 1
	if i > last {
		i = last
	}
	r := t.In(i)
	if t.IsVariadic() && i == last {
		r = r.Elem()
	}
	return r
}

func reflectTypeOf(i interface{}) reflectType {
	return reflectType{reflect.TypeOf(i)}
}

// Call calls function fn with given parameters.
func Call(fn interface{}, params ...interface{}) []interface{} {
	fnR := reflect.ValueOf(fn)

	if reflect.ValueOf(fn).Kind() != reflect.Func {
		panic(fmt.Sprintf("func Call expects func, %T given", fn))
	}

	fnType := reflectTypeOf(fn)

	if len(params) > fnType.NumIn() && !fnType.IsVariadic() {
		panic("Call with too many input arguments")
	}

	paramsRef := make([]reflect.Value, len(params))
	for i, p := range params {
		paramsRef[i] = reflect.ValueOf(p).Convert(fnType.inVariadicAware(i))
	}

	var result []interface{}
	for _, v := range fnR.Call(paramsRef) {
		result = append(result, v.Interface())
	}

	return result
}

// SafeCall calls function Call and returns error in case of panic.
func SafeCall(fn interface{}, params ...interface{}) (result []interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%s", p)
		}
	}()

	return Call(fn, params...), nil
}

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

// CallProvider calls function given as first param with given parameters,
// function must return 1 or 2 values, second value must be error if exists.
func CallProvider(provider interface{}, params ...interface{}) (interface{}, error) {
	t := reflect.TypeOf(provider)
	if t.NumOut() == 0 || t.NumOut() > 2 {
		return nil, fmt.Errorf("provider must return 1 or 2 values, given function returns %d values", t.NumOut())
	}
	if t.NumOut() == 2 && !t.Out(1).Implements(errorInterface) {
		return nil, fmt.Errorf("second value returned by provider must implements error interface")
	}

	results, err := SafeCall(provider, params...)
	if err != nil {
		return nil, err
	}

	r := results[0]
	var e error
	if len(results) > 1 {
		// do not panic when results[1] == nil
		e, _ = results[1].(error)
	}

	return r, e
}
