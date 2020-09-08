package caller

import (
	"fmt"
	"reflect"

	helpersReflect "github.com/gomponents/gontainer-helpers/reflect"
)

func call(fn reflect.Value, params ...interface{}) []interface{} {
	if fn.Kind() != reflect.Func {
		panic(fmt.Sprintf("expects %s, %T given", reflect.Func.String(), fn.Type().String()))
	}

	fnType := reflectType{fn.Type()}

	if len(params) > fnType.NumIn() && !fnType.IsVariadic() {
		panic("call with too many input arguments")
	}

	minParams := fnType.NumIn()
	if fnType.IsVariadic() {
		minParams--
	}
	if len(params) < minParams {
		panic("call with too few input arguments")
	}

	paramsRef := make([]reflect.Value, len(params))
	for i, p := range params {
		vp := reflect.ValueOf(p)
		convertTo := fnType.inVariadicAware(i)
		cp, ok := helpersReflect.Convert(vp, convertTo)
		if !ok {
			panic(fmt.Sprintf("arg%d: cannot cast `%s` to `%s`", i, vp.Type().String(), convertTo.String()))
		}
		paramsRef[i] = cp
	}

	var result []interface{}
	for _, v := range fn.Call(paramsRef) {
		result = append(result, v.Interface())
	}

	return result
}

// MustCall calls function fn with given parameters. It panics in case of error, use Call to avoid panic.
func MustCall(fn interface{}, params ...interface{}) []interface{} {
	return call(reflect.ValueOf(fn), params...)
}

// Call calls function MustCall and returns error in case of panic.
func Call(fn interface{}, params ...interface{}) (result []interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%s", p)
		}
	}()

	return MustCall(fn, params...), nil
}

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

// CallProvider calls function given as first param with given parameters,
// function must return 1 or 2 values, second value must be error if exists.
func CallProvider(provider interface{}, params ...interface{}) (interface{}, error) {
	t := reflect.TypeOf(provider)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("provider must be kind of %s, %s given", reflect.Func.String(), t.Kind().String())
	}
	if t.NumOut() == 0 || t.NumOut() > 2 {
		return nil, fmt.Errorf("provider must return 1 or 2 values, given function returns %d values", t.NumOut())
	}
	if t.NumOut() == 2 && !t.Out(1).Implements(errorInterface) {
		return nil, fmt.Errorf("second value returned by provider must implements error interface")
	}

	results, err := Call(provider, params...)
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

func MustCallByName(object interface{}, method string, params ...interface{}) []interface{} {
	val := reflect.ValueOf(object)
	fn := val.MethodByName(method)
	for !fn.IsValid() && (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) {
		val = val.Elem()
		fn = val.MethodByName(method)
	}

	if !fn.IsValid() {
		panic(fmt.Sprintf("invalid func `%T`.`%s`", object, method))
	}

	return call(fn, params...)
}

func CallByName(object interface{}, method string, params ...interface{}) (result []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	return MustCallByName(object, method, params...), nil
}

func MustCallWitherByName(object interface{}, wither string, params ...interface{}) interface{} {
	val := reflect.ValueOf(object)
	fn := val.MethodByName(wither)
	for !fn.IsValid() && (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) {
		val = val.Elem()
		fn = val.MethodByName(wither)
	}

	if !fn.IsValid() {
		panic(fmt.Sprintf("invalid wither `%T`.`%s`", object, wither))
	}

	t := fn.Type()

	if t.NumOut() != 1 {
		panic(fmt.Sprintf("wither must return 1 value, given function returns %d values", t.NumOut()))
	}

	return call(fn, params...)[0]
}

func CallWitherByName(object interface{}, wither string, params ...interface{}) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	return MustCallWitherByName(object, wither, params...), nil
}

func CallDecorator(decorator interface{}, object interface{}, args ...interface{}) (interface{}, error) {
	t := reflect.TypeOf(decorator)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("decorator must be kind of %s, %s given", reflect.Func.String(), t.Kind().String())
	}
	if t.NumOut() == 0 || t.NumOut() > 2 {
		return nil, fmt.Errorf("decorator must return 1 or 2 values, given function returns %d values", t.NumOut())
	}
	if t.NumOut() == 2 && !t.Out(1).Implements(errorInterface) {
		return nil, fmt.Errorf("second value returned by provider must implements error interface")
	}

	results, err := Call(decorator, append([]interface{}{object}, args...))
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

type reflectType struct {
	reflect.Type
}

// inVariadicAware works almost same as reflect.Type.In,
// but it returns t.In(t.NumIn() - 1).Elem() for t.isVariadic() && i >= t.NumIn().
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
