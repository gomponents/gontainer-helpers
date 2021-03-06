package exporters

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myString string
type aliasString = string
type myInt int
type aliasInt = int
type myBool bool
type aliasBool = bool

type mockExporter struct {
	result string
	error  error
}

func (m mockExporter) Export(interface{}) (string, error) {
	return m.result, m.error
}

func TestChainExporter_Export(t *testing.T) {
	exporter := NewDefaultExporter()

	t.Run("Given scenarios", func(t *testing.T) {
		scenarios := map[string]struct {
			input  interface{}
			output string
			error  string
		}{
			"nil": {
				input:  nil,
				output: "nil",
			},
			"false": {
				input:  false,
				output: "false",
			},
			"true": {
				input:  true,
				output: "true",
			},
			"123": {
				input:  int(123),
				output: "int(123)",
			},
			"`hello world`": {
				input:  "hello world",
				output: `"hello world"`,
			},
			"[]byte": {
				input:  []byte("hello world 你好，世界"),
				output: `[]byte("hello world \u4f60\u597d\uff0c\u4e16\u754c")`,
			},
			"struct {}": {
				input: struct{}{},
				error: "parameter of type `struct {}` is not supported",
			},
			"*testing.T": {
				input: t,
				error: "parameter of type `*testing.T` is not supported",
			},
			`myString("foo")`: {
				input: myString("foo"),
				error: "parameter of type `exporters.myString` is not supported",
			},
			`aliasString("foo")`: {
				input:  aliasString("foo"),
				output: `"foo"`,
			},
			`myInt(5)`: {
				input: myInt(5),
				error: "parameter of type `exporters.myInt` is not supported",
			},
			`aliasInt(5)`: {
				input:  aliasInt(5),
				output: "int(5)",
			},
			`myBool(true)`: {
				input: myBool(true),
				error: "parameter of type `exporters.myBool` is not supported",
			},
			`aliasBool(true)`: {
				input:  aliasBool(true),
				output: "true",
			},
		}

		for k, s := range scenarios {
			t.Run(k, func(t *testing.T) {
				output, err := exporter.Export(s.input)
				if s.error != "" {
					assert.EqualError(t, err, s.error)
					assert.Equal(t, "", output)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, s.output, output)
			})
		}
	})
}

func TestExport(t *testing.T) {
	scenarios := []struct {
		input  interface{}
		output string
		error  string
		panic  string
	}{
		{
			input:  123,
			output: "int(123)",
		},
		{
			input:  []interface{}{1, "2", 3.14},
			output: `[]interface{}{int(1), "2", float64(3.14)}`,
		},
		{
			input:  [3]interface{}{1, "2", 3.14},
			output: `[3]interface{}{int(1), "2", float64(3.14)}`,
		},
		{
			input:  []interface{}{},
			output: "make([]interface{}, 0)",
		},
		{
			input:  [0]interface{}{},
			output: "[0]interface{}{}",
		},
		{
			input: []interface{}{struct{}{}},
			error: "cannot export elem 0 of slice: parameter of type `struct {}` is not supported",
			panic: "cannot export parameter of type `[]interface {}` to string: cannot export elem 0 of slice: parameter of type `struct {}` is not supported",
		},
		{
			input:  []int{1, 2, 3},
			output: "[]int{int(1), int(2), int(3)}",
		},
		{
			input:  [3]int{1, 2, 3},
			output: "[3]int{int(1), int(2), int(3)}",
		},
		{
			input:  []float32{},
			output: "make([]float32, 0)",
		},
		{
			input:  [0]float32{},
			output: "[0]float32{}",
		},
		{
			input: struct{}{},
			error: "parameter of type `struct {}` is not supported",
			panic: "cannot export parameter of type `struct {}` to string: parameter of type `struct {}` is not supported",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			func() {
				defer func() {
					r := recover()
					if s.panic == "" {
						assert.Nil(t, r)
						return
					}
					assert.Equal(t, s.panic, r)
				}()
				assert.Equal(t, s.output, MustExport(s.input))
			}()

			o, err := Export(s.input)
			if s.error == "" {
				assert.NoError(t, err)
				assert.Equal(t, s.output, o)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}

	t.Run("Given invalid scenario", func(t *testing.T) {
		originalExporter := defaultExporter
		defer func() {
			defaultExporter = originalExporter
		}()

		expectedErr := fmt.Errorf("my test error")
		defaultExporter = mockExporter{
			error: expectedErr,
		}
		_, err := Export(123)
		assert.EqualError(t, err, expectedErr.Error())
	})
}

func TestToString(t *testing.T) {
	scenarios := []struct {
		input  interface{}
		output string
		error  string
	}{
		{
			input:  true,
			output: "true",
		},
		{
			input:  nil,
			output: "nil",
		},
		{
			input: struct{}{},
			error: "parameter of type `struct {}` is not supported",
		},
		{
			input:  "Mary Jane",
			output: "Mary Jane",
		},
		{
			input:  int(5),
			output: "5",
		},
		{
			input:  float64(3.14),
			output: "3.14",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			t.Run("ToString", func(t *testing.T) {
				result, err := ToString(s.input)

				if s.error != "" {
					assert.Empty(t, result)
					assert.EqualError(t, err, s.error)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, s.output, result)
			})

			t.Run("MustToString", func(t *testing.T) {
				defer func() {
					err := recover()
					if s.error == "" {
						assert.Nil(t, err)
						return
					}

					assert.NotNil(t, err)

					assert.Equal(
						t,
						fmt.Sprintf(
							"cannot cast parameter of type `%T` to string: %s",
							s.input,
							s.error,
						),
						fmt.Sprintf("%s", err),
					)
				}()

				assert.Equal(t, s.output, MustToString(s.input))
			})
		})
	}
}

func TestNumericExporter_Supports(t *testing.T) {
	scenarios := []struct {
		input    interface{}
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    math.Pi,
			expected: true,
		},
		{
			input:    "3.14",
			expected: false,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(
				t,
				s.expected,
				NumericExporter{}.Supports(s.input),
			)
		})
	}
}

func TestPrimitiveTypeSliceExporter_Supports(t *testing.T) {
	scenarios := []struct {
		input    interface{}
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    math.Pi,
			expected: false,
		},
		{
			input:    "3.14",
			expected: false,
		},
		{
			input:    []uint{0, 1},
			expected: true,
		},
		{
			input:    []struct{}{},
			expected: false,
		},
		{
			input:    []myBool{},
			expected: false,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(
				t,
				s.expected,
				PrimitiveTypeSliceExporter{}.Supports(s.input),
			)
		})
	}
}

func TestInterfaceSliceExporter_Supports(t *testing.T) {
	scenarios := []struct {
		input    interface{}
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    math.Pi,
			expected: false,
		},
		{
			input:    "3.14",
			expected: false,
		},
		{
			input:    []uint{0, 1},
			expected: false,
		},
		{
			input:    []struct{}{},
			expected: false,
		},
		{
			input:    []interface{}{},
			expected: true,
		},
		{
			input:    []myBool{},
			expected: false,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			assert.Equal(
				t,
				s.expected,
				InterfaceSliceExporter{}.Supports(s.input),
				fmt.Sprintf("value: %#v", s.input),
			)
		})
	}
}

func TestPrimitiveTypeSliceExporter_Export(t *testing.T) {
	t.Run("Given error in subexporter", func(t *testing.T) {
		exp := PrimitiveTypeSliceExporter{
			exporter: ChainExporter{},
		}
		v, err := exp.Export([]uint{1})
		assert.Equal(t, "", v)
		assert.EqualError(
			t,
			err,
			"unexpected err when exporting elem 0: parameter of type `uint` is not supported",
		)
	})
}
