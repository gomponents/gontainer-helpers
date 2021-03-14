// +build limitations

package caller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type book struct {
	title string
}

func (b *book) SetTitle(t string) {
	b.title = t
}

func (b book) WithTitle(t string) book {
	b.title = t
	return b
}

func TestLimitations(t *testing.T) {
	// Method with pointer receiver requires explicit definition of pointer:
	// v := &book{}; CallByName(v, ...
	// var v interface{} = &book{}; CallByName(v, ...
	// v := book{}; CallByName(&v, ...
	//
	// Creating variable as a value will not work:
	// v := book{}; CallByName(v, ...
	// var v interface = book{}; CallByName(&v, ...
	t.Run("Call method", func(t *testing.T) {
		t.Run("Pointer receiver", func(t *testing.T) {
			t.Run("Given errors", func(t *testing.T) {
				t.Run("v := book{}; CallByName(v, ...", func(t *testing.T) {
					harryPotter := book{}
					r, err := CallByName(harryPotter, "SetTitle", "Harry Potter")
					assert.EqualError(t, err, "invalid func `caller.book`.`SetTitle`")
					assert.Nil(t, r)
					assert.Zero(t, harryPotter)
				})
				t.Run("var v interface{} = book{}; CallByName(&v, ...", func(t *testing.T) {
					var harryPotter interface{} = book{}
					r, err := CallByName(&harryPotter, "SetTitle", "Harry Potter")
					assert.EqualError(t, err, "invalid func `*interface {}`.`SetTitle`")
					assert.Nil(t, r)
					assert.Equal(t, book{}, harryPotter)
				})
			})
			t.Run("Given scenarios", func(t *testing.T) {
				t.Run("v := book{}; CallByName(&v, ...", func(t *testing.T) {
					harryPotter := book{}
					r, err := CallByName(&harryPotter, "SetTitle", "Harry Potter")
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, "Harry Potter", harryPotter.title)
				})
				t.Run("v := &book{}; CallByName(&v, ...", func(t *testing.T) {
					harryPotter := &book{}
					r, err := CallByName(&harryPotter, "SetTitle", "Harry Potter")
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, "Harry Potter", harryPotter.title)
				})
				t.Run("v := &book{}", func(t *testing.T) {
					harryPotter := &book{}
					r, err := CallByName(harryPotter, "SetTitle", "Harry Potter")
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, "Harry Potter", harryPotter.title)
				})
				t.Run("var v interface{} = &book{}; CallByName(v, ...", func(t *testing.T) {
					var harryPotter interface{} = &book{}
					r, err := CallByName(harryPotter, "SetTitle", "Harry Potter")
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &book{title: "Harry Potter"}, harryPotter)
				})
				t.Run("var v interface{} = &book{}; CallByName(&v, ...", func(t *testing.T) {
					var harryPotter interface{} = &book{}
					r, err := CallByName(&harryPotter, "SetTitle", "Harry Potter")
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &book{title: "Harry Potter"}, harryPotter)
				})
				t.Run("var v interface{ SetTitle(string) } = &book{}; CallByName(v, ...", func(t *testing.T) {
					var harryPotter interface{ SetTitle(string) } = &book{}
					r, err := CallByName(harryPotter, "SetTitle", "Harry Potter")
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &book{title: "Harry Potter"}, harryPotter)
				})
			})
		})
		// Methods with value receiver do not have any limitations
		t.Run("Value receiver", func(t *testing.T) {
			t.Run("b := book{}", func(t *testing.T) {
				b := book{}
				r, err := CallWitherByName(b, "WithTitle", "Harry Potter")
				assert.NoError(t, err)
				assert.Equal(t, book{title: "Harry Potter"}, r)
				assert.Zero(t, b)
			})
			t.Run("b := &book{}", func(t *testing.T) {
				b := book{}
				r, err := CallWitherByName(b, "WithTitle", "Harry Potter")
				assert.NoError(t, err)
				assert.Equal(t, book{title: "Harry Potter"}, r)
				assert.Zero(t, b)
			})
			t.Run("var b interface{} = book{}", func(t *testing.T) {
				var b interface{} = book{}
				r, err := CallWitherByName(b, "WithTitle", "Harry Potter")
				assert.NoError(t, err)
				assert.Equal(t, book{title: "Harry Potter"}, r)
				assert.Zero(t, b)
			})
			t.Run("var b interface{} = &book{}", func(t *testing.T) {
				var b interface{} = &book{}
				r, err := CallWitherByName(b, "WithTitle", "Harry Potter")
				assert.NoError(t, err)
				assert.Equal(t, book{title: "Harry Potter"}, r)
				assert.Equal(t, &book{}, b)
			})
		})
	})
}
