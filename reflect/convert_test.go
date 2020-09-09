package reflect

import (
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	// todo
	t.Run("Convert parameters", func(t *testing.T) {
		//float64Val := float64(5)

		scenarios := map[string]struct {
			input  interface{}
			output interface{}
			ok     bool
		}{
			"[]interface{} to []type (correct)": {
				input:  []interface{}{1, 2, 3},
				output: []int{1, 2, 3},
				ok:     true,
			},
			"[]interface{} to []type (incorrect)": {
				input:  []struct{}{},
				output: []int{},
				ok:     false,
			},
			"float64 to int": {
				input:  float64(math.Pi),
				output: 3,
				ok:     true,
			},
			//"*float64 to *int": {
			//	fn:    func(*int) {},
			//	input: &float64Val,
			//	error: "arg0: cannot cast `*float64` to `*int`",
			//},
			//"*float64 to *float32": {
			//	fn:    func(*float32) {},
			//	input: &float64Val,
			//	error: "arg0: cannot cast `*float64` to `*float32`",
			//},
			//"int to float64": {
			//	fn: func(v float64) float64 {
			//		return v
			//	},
			//	input:  int(5),
			//	output: float64(5),
			//},
			//"string to []byte": {
			//	fn: func(v []byte) []byte {
			//		return v
			//	},
			//	input:  "hello",
			//	output: []byte("hello"),
			//},
			//"[]byte to string": {
			//	fn: func(v string) string {
			//		return v
			//	},
			//	input:  []byte("hello"),
			//	output: "hello",
			//},
			//"string to int": { // cannot convert string to int
			//	fn:    func(int) {},
			//	input: "5",
			//	error: "arg0: cannot cast `string` to `int`",
			//},
			//"int to string": { // but reverse conversion is possible, isn't worth to unify this behavior?
			//	fn: func(v string) string {
			//		return v
			//	},
			//	input:  5,
			//	output: "\x05",
			//},
			//"zero value": {
			//	fn: func(v int) int {
			//		return v
			//	},
			//	input:  nil,
			//	output: 0,
			//},
		}

		for n, s := range scenarios {
			t.Run(n, func(t *testing.T) {
				v, ok := Convert(reflect.ValueOf(s.input), reflect.TypeOf(s.output))
				assert.Equal(t, s.ok, ok)
				if ok {
					assert.Equal(t, s.output, v.Interface())
				}
			})
		}
	})
}
