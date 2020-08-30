package exporters

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	defaultExporter     = NewDefaultExporter()
	defaultStringCaster = NewChainExporter(
		&BoolExporter{},
		&NilExporter{},
		&NumericExporter{},
	)
)

// Export exports input value to GO code
func Export(i interface{}) (string, error) {
	return defaultExporter.Export(i)
}

// ToString casts input value to string.
// In case of string value output will be equal to input in opposition to Export function.
func ToString(i interface{}) (string, error) {
	if r, ok := i.(string); ok {
		return r, nil
	}

	return defaultStringCaster.Export(i)
}

// MustToString casts input value to string
func MustToString(i interface{}) string {
	r, err := ToString(i)
	if err != nil {
		panic(fmt.Sprintf("cannot cast parameter of type `%T` to string: %s", i, err.Error()))
	}
	return r
}

type Exporter interface {
	Export(interface{}) (string, error)
}

type SubExporter interface {
	Exporter
	Supports(interface{}) bool
}

type ChainExporter struct {
	exporters []SubExporter
}

func (c ChainExporter) Export(v interface{}) (string, error) {
	for _, e := range c.exporters {
		if e.Supports(v) {
			return e.Export(v)
		}
	}

	return "", errors.New(fmt.Sprintf("parameter of type `%T` is not supported", v))
}

func NewDefaultExporter() Exporter {
	interfaceSliceExporter := NewInterfaceSliceExporter(nil)
	primitiveTypeSliceExporter := NewPrimitiveTypeSliceExporter(nil)

	result := NewChainExporter(
		&BoolExporter{},
		&NilExporter{},
		&NumericExporter{},
		&StringExporter{},
		interfaceSliceExporter,
		primitiveTypeSliceExporter,
	)
	interfaceSliceExporter.exporter = result
	primitiveTypeSliceExporter.exporter = result

	return result
}

func NewChainExporter(exporters ...SubExporter) *ChainExporter {
	return &ChainExporter{exporters: exporters}
}

type BoolExporter struct{}

func (b BoolExporter) Export(v interface{}) (string, error) {
	if v == true {
		return "true", nil
	}

	return "false", nil
}

func (b BoolExporter) Supports(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

type NilExporter struct{}

func (n NilExporter) Export(interface{}) (string, error) {
	return "nil", nil
}

func (n NilExporter) Supports(v interface{}) bool {
	return v == nil
}

type NumericExporter struct{}

func (n NumericExporter) Export(v interface{}) (string, error) {
	switch reflect.TypeOf(v).Kind() {
	case
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return fmt.Sprintf("%#v", v), nil
	}
	return fmt.Sprintf("%d", v), nil
}

var (
	numericKinds = []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
	}
)

func (n NumericExporter) Supports(v interface{}) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}

	for _, k := range numericKinds {
		if k == t.Kind() {
			return true
		}
	}

	return false
}

type StringExporter struct{}

func (s StringExporter) Export(v interface{}) (string, error) {
	return fmt.Sprintf("%+q", v), nil
}

func (s StringExporter) Supports(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

type InterfaceSliceExporter struct {
	exporter Exporter
}

func NewInterfaceSliceExporter(exporter Exporter) *InterfaceSliceExporter {
	return &InterfaceSliceExporter{exporter: exporter}
}

func (i InterfaceSliceExporter) Export(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	if val.Len() == 0 {
		return "make([]interface{}, 0)", nil
	}
	parts := make([]string, val.Len())
	for j := 0; j < val.Len(); j++ {
		part, err := i.exporter.Export(val.Index(j).Interface())
		if err != nil {
			return "", fmt.Errorf("cannot export elem %d of slice: %s", j, err.Error())
		}
		parts[j] = part
	}

	return "[]interface{}{" + strings.Join(parts, ", ") + "}", nil
}

func (i InterfaceSliceExporter) Supports(v interface{}) bool {
	val := reflect.ValueOf(v)
	return val.Type().Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Interface
}

type PrimitiveTypeSliceExporter struct {
	exporter Exporter
}

func NewPrimitiveTypeSliceExporter(exporter Exporter) *PrimitiveTypeSliceExporter {
	return &PrimitiveTypeSliceExporter{exporter: exporter}
}

func (p PrimitiveTypeSliceExporter) Export(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	if val.Len() == 0 {
		return fmt.Sprintf("make([]%s, 0)", val.Type().Elem().Kind().String()), nil
	}
	parts := make([]string, val.Len())
	for i := 0; i < val.Len(); i++ {
		var err error
		parts[i], err = p.exporter.Export(val.Index(i).Interface())
		if err != nil {
			return "", fmt.Errorf("unexpected err when exporting elem %d: %s", i, err.Error())
		}
	}
	return "[]" + val.Type().Elem().Kind().String() + "{" + strings.Join(parts, ", ") + "}", nil
}

func (p PrimitiveTypeSliceExporter) Supports(v interface{}) bool {
	val := reflect.ValueOf(v)
	if val.Type().Kind() != reflect.Slice {
		return false
	}

	switch val.Type().Elem().Kind() {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return true
	}

	return false
}
