package reflect

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Run("Convert parameters", func(t *testing.T) {
		float64Val := float64(5)

		scenarios := map[string]struct {
			input  interface{}
			output interface{}
			error  string
		}{
			"[][]interface{} to [][]int": {
				input:  [][]interface{}{{1, 2, 3}},
				output: [][]int{{1, 2, 3}},
			},
			"[][]interface{} to [][]int (invalid)": {
				input:  [][]interface{}{{1, false, 3}},
				output: [][]int{{1, 2, 3}},
				error:  "cannot cast `[][]interface {}` to `[][]int`: el0: cannot cast `[]interface {}` to `[]int`: el1: cannot cast `bool` to `int`",
			},
			"[][]int to [][]interface{}": {
				input:  [][]int{{1, 2, 3}},
				output: [][]interface{}{{1, 2, 3}},
			},
			"[][]uint to [][]int": {
				input:  [][]uint{{1, 2, 3}},
				output: [][]int{{1, 2, 3}},
			},
			"[]interface{} to []int": {
				input:  []interface{}{1, 2, 3},
				output: []int{1, 2, 3},
			},
			"[]interface{} to []int (invalid #1)": {
				input:  []interface{}{1, 2, nil},
				output: []int{},
				error:  "cannot cast `[]interface {}` to `[]int`: el2: cannot cast `<nil>` to `int`",
			},
			"[]interface{} to []int (invalid #2)": {
				input:  []interface{}{1, 2, 3, struct{}{}},
				output: []int{},
				error:  "cannot cast `[]interface {}` to `[]int`: el3: cannot cast `struct {}` to `int`",
			},
			"[]interface{} to []*int": {
				input:  []interface{}{nil, nil},
				output: []*int{nil, nil},
			},
			"[]int to []interface{}": {
				input:  []int{1, 2, 3},
				output: []interface{}{1, 2, 3},
			},
			"[]interface{} to []interface{}": {
				input:  []interface{}{1, 2, 3},
				output: []interface{}{1, 2, 3},
			},
			"[]int8 to []int": {
				input:  []int8{1, 2, 3},
				output: []int{1, 2, 3},
			},
			"[]int to []int8": {
				input:  []int{1, 2, 256},
				output: []int8{1, 2, 0},
			},
			"[]struct{}{} to []type": {
				input:  []struct{}{},
				output: []int{},
				error:  "cannot cast `[]struct {}` to `[]int`",
			},
			"float64 to int": {
				input:  float64(math.Pi),
				output: 3,
			},
			"nil to int": {
				input:  nil,
				output: 0,
				error:  "cannot cast `<nil>` to `int`",
			},
			"nil to *int": {
				input:  nil,
				output: (*int)(nil),
			},
			"*float64 to *int": {
				input:  &float64Val,
				output: (*int)(nil),
				error:  "cannot cast `*float64` to `*int`",
			},
			"*float64 to *float32": {
				input:  &float64Val,
				output: (*float32)(nil),
				error:  "cannot cast `*float64` to `*float32`",
			},
			"*float64 to *float64": {
				input:  &float64Val,
				output: &float64Val,
			},
			"int to float64": {
				input:  int(5),
				output: float64(5),
			},
			"string to []byte": {
				input:  "hello",
				output: []byte("hello"),
			},
			"[]byte to string": {
				input:  []byte("hello"),
				output: "hello",
			},
			"string to int": { // cannot convert string to int
				input:  "5",
				output: int(5),
				error:  "cannot cast `string` to `int`",
			},
			"int to string": { // but reverse conversion is possible, isn't worth to unify this behavior?
				input:  5,
				output: "\x05",
			},
			"nil to *interface{}": {
				input:  nil,
				output: (*interface{})(nil),
			},
		}

		for n, s := range scenarios {
			t.Run(n, func(t *testing.T) {
				v, err := Convert(s.input, reflect.TypeOf(s.output))
				if s.error != "" {
					assert.EqualError(t, err, s.error)
					return
				}
				if assert.NoError(t, err) {
					assert.Equal(t, s.output, v.Interface())
				}
			})
		}
	})
}
