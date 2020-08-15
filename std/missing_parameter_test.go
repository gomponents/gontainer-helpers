package std

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMissingParameter(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		defer func() {
			assert.Equal(t, recover(), "missing parameter")
		}()
		GetMissingParameter()
	})
	t.Run("With argument", func(t *testing.T) {
		msg := fmt.Sprintf("%f", rand.Float32())
		defer func() {
			assert.Equal(t, recover(), msg)
		}()
		GetMissingParameter(msg)
	})
}
