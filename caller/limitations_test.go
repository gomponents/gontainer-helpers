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
	const harryPotterTitle = "Harry Potter"
	harryPotter := book{title: harryPotterTitle}
	emptyBook := book{}

	// https://github.com/golang/go/wiki/MethodSets#interfaces

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
					b := book{}
					r, err := CallByName(b, "SetTitle", harryPotterTitle)
					assert.EqualError(t, err, "invalid func `caller.book`.`SetTitle`")
					assert.Nil(t, r)
					assert.Zero(t, b)
				})
				t.Run("var v interface{} = book{}; CallByName(&v, ...", func(t *testing.T) {
					var b interface{} = book{}
					r, err := CallByName(&b, "SetTitle", harryPotterTitle)
					assert.EqualError(t, err, "invalid func `*interface {}`.`SetTitle`")
					assert.Nil(t, r)
					assert.Equal(t, emptyBook, b)
				})
			})
			t.Run("Given scenarios", func(t *testing.T) {
				t.Run("v := book{}; CallByName(&v, ...", func(t *testing.T) {
					b := book{}
					r, err := CallByName(&b, "SetTitle", harryPotterTitle)
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, harryPotter, b)
				})
				t.Run("v := &book{}; CallByName(&v, ...", func(t *testing.T) {
					b := &book{}
					r, err := CallByName(&b, "SetTitle", harryPotterTitle)
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("v := &book{}; CallByName(v, ...", func(t *testing.T) {
					b := &book{}
					r, err := CallByName(b, "SetTitle", harryPotterTitle)
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("var v interface{} = &book{}; CallByName(v, ...", func(t *testing.T) {
					var b interface{} = &book{}
					r, err := CallByName(b, "SetTitle", harryPotterTitle)
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("var v interface{} = &book{}; CallByName(&v, ...", func(t *testing.T) {
					var b interface{} = &book{}
					r, err := CallByName(&b, "SetTitle", harryPotterTitle)
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &harryPotter, b)
				})
				t.Run("var v interface{ SetTitle(string) } = &book{}; CallByName(v, ...", func(t *testing.T) {
					var b interface{ SetTitle(string) } = &book{}
					r, err := CallByName(b, "SetTitle", harryPotterTitle)
					assert.NoError(t, err)
					assert.Nil(t, r)
					assert.Equal(t, &harryPotter, b)
				})
			})
		})
		// Methods with value receiver do not have any limitations
		t.Run("Value receiver", func(t *testing.T) {
			t.Run("b := book{}", func(t *testing.T) {
				b := book{}
				r, err := CallWitherByName(b, "WithTitle", harryPotterTitle)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Zero(t, b)
			})
			t.Run("b := &book{}", func(t *testing.T) {
				b := &book{}
				r, err := CallWitherByName(b, "WithTitle", harryPotterTitle)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Equal(t, &emptyBook, b)
			})
			t.Run("var b interface{} = book{}", func(t *testing.T) {
				var b interface{} = book{}
				r, err := CallWitherByName(b, "WithTitle", harryPotterTitle)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Equal(t, emptyBook, b)
			})
			t.Run("var b interface{} = &book{}", func(t *testing.T) {
				var b interface{} = &book{}
				r, err := CallWitherByName(b, "WithTitle", harryPotterTitle)
				assert.NoError(t, err)
				assert.Equal(t, harryPotter, r)
				assert.Equal(t, &emptyBook, b)
			})
		})
	})
}
