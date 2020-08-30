package exporters

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainExporter_Export(t *testing.T) {
	exporter := NewDefaultExporter()

	t.Run("Given valid scenarios", func(t *testing.T) {
		scenarios := map[string]struct {
			input  interface{}
			output string
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
				input:  123,
				output: "123",
			},
			"`hello world`": {
				input:  "hello world",
				output: `"hello world"`,
			},
			"complex64(0.588)": {
				input:  complex64(0.588),
				output: "(0.588+0i)",
			},
			"complex128(0.588)": {
				input:  complex128(0.588),
				output: "(0.588+0i)",
			},
			"complex128(3.14)": {
				input:  complex128(3.14),
				output: "(3.14+0i)",
			},
		}

		for k, s := range scenarios {
			t.Run(k, func(t *testing.T) {
				output, err := exporter.Export(s.input)
				assert.NoError(t, err)
				assert.Equal(t, s.output, output)
			})
		}
	})

	t.Run("Given invalid scenarios", func(t *testing.T) {
		scenarios := map[string]struct {
			input interface{}
			error string
		}{
			"struct {}": {
				input: struct{}{},
				error: "parameter of type `struct {}` is not supported",
			},
			"*testing.T": {
				input: t,
				error: "parameter of type `*testing.T` is not supported",
			},
		}

		for k, s := range scenarios {
			t.Run(k, func(t *testing.T) {
				output, err := exporter.Export(s.input)
				assert.EqualError(t, err, s.error)
				assert.Equal(t, "", output)
			})
		}
	})
}

func TestExport(t *testing.T) {
	scenarios := []struct {
		input  interface{}
		output string
		error  string
	}{
		{
			input:  123,
			output: "123",
		},
		{
			input:  []interface{}{1, "2", 3.14},
			output: `[]interface{}{1, "2", 3.14}`,
		},
		{
			input:  [3]interface{}{1, "2", 3.14},
			output: `[3]interface{}{1, "2", 3.14}`,
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
		},
		{
			input:  []int{1, 2, 3},
			output: "[]int{1, 2, 3}",
		},
		{
			input:  [3]int{1, 2, 3},
			output: "[3]int{1, 2, 3}",
		},
		{
			input:  []float32{},
			output: "make([]float32, 0)",
		},
		{
			input:  [0]float32{},
			output: "[0]float32{}",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
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
		{
			input:    complex128(1),
			expected: true,
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

type mockExporter struct {
	result string
	error  error
}

func (m mockExporter) Export(interface{}) (string, error) {
	return m.result, m.error
}
