package exporters

import (
	"fmt"
	"math"
	"math/rand"
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
	originalExporter := defaultExporter
	defer func() {
		defaultExporter = originalExporter
	}()

	t.Run("Given valid scenario", func(t *testing.T) {
		expected := fmt.Sprintf("%f", rand.Float32())
		defaultExporter = mockExporter{
			result: expected,
		}
		result, err := Export(123)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Given invalid scenario", func(t *testing.T) {
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
