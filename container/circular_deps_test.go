package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newCircularDeps(t *testing.T) {
	t.Run("No circular deps", func(t *testing.T) {
		d := newCircularDeps()

		assert.Nil(t, d.start("foo"))
		assert.Nil(t, d.start("bar"))
		d.stop()
		d.stop()

		assert.Empty(t, d.chain)
	})

	t.Run("Circular deps", func(t *testing.T) {
		d := newCircularDeps()

		assert.Nil(t, d.start("app"))
		assert.Nil(t, d.start("storage"))
		assert.Nil(t, d.start("db"))
		assert.Equal(t, []string{"app", "storage", "db", "storage"}, d.start("storage"))
		d.stop()
		d.stop()
		d.stop()
		d.stop()

		assert.Empty(t, d.chain)
	})
}
