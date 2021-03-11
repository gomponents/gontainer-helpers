package reflect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	t.Run("var from interface{}", func(t *testing.T) {
		var from interface{} = car{age: 5}
		var to car

		assert.NoError(t, Copy(from, &to))
		assert.Equal(t, 5, to.age)

		t.Run("MustCopy", func(t *testing.T) {
			to.age = 0
			MustCopy(from, &to)
			assert.Equal(t, 5, to.age)
		})
	})
	t.Run("var to interface{}", func(t *testing.T) {
		six := 6
		scenarios := []interface{}{
			5,
			3.14,
			struct{}{},
			nil,
			&six,
			car{age: 10},
			&car{age: 10},
			(*car)(nil),
		}

		for id, d := range scenarios {
			t.Run(fmt.Sprintf("%d: %T", id, d), func(t *testing.T) {
				var to interface{}
				assert.NoError(t, Copy(d, &to))
				assert.Equal(t, d, to)
				if reflect.ValueOf(d).Kind() == reflect.Ptr {
					assert.Same(t, d, to)
				}

				t.Run("MustCopy", func(t *testing.T) {
					var to interface{}
					MustCopy(d, &to)
					assert.Equal(t, d, to)
					if reflect.ValueOf(d).Kind() == reflect.Ptr {
						assert.Same(t, d, to)
					}
				})
			})
		}
	})
	t.Run("Given errors", func(t *testing.T) {
		t.Run("non-pointer value", func(t *testing.T) {
			const msg = "reflect.Copy expects pointer in second argument, int given"
			assert.EqualError(
				t,
				Copy(5, 5),
				msg,
			)

			t.Run("MustCopy", func(t *testing.T) {
				defer func() {
					assert.Equal(
						t,
						msg,
						fmt.Sprintf("%s", recover()),
					)
				}()

				MustCopy(5, 5)
			})
		})
	})
}

type car struct {
	age int
}
