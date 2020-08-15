package std

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustGetMissingParameter(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		defer func() {
			assert.Equal(t, recover(), "missing parameter")
		}()
		MustGetMissingParameter()
	})
	t.Run("With argument", func(t *testing.T) {
		msg := fmt.Sprintf("%f", rand.Float32())
		defer func() {
			assert.Equal(t, recover(), msg)
		}()
		MustGetMissingParameter(msg)
	})
}
