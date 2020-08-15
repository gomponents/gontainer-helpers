package exporters

import (
	"errors"
	"fmt"
	"reflect"
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
	return NewChainExporter(
		&BoolExporter{},
		&NilExporter{},
		&NumericExporter{},
		&StringExporter{},
	)
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
