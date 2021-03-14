package std

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterTodo(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		v, err := ParameterTodo()
		assert.EqualError(t, err, paramTodoMsg)
		assert.Nil(t, v)
	})
	t.Run("With argument", func(t *testing.T) {
		msg := fmt.Sprintf("%f", rand.Float32())
		v, err := ParameterTodo(msg)
		assert.EqualError(t, err, msg)
		assert.Nil(t, v)
	})
}
