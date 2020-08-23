package setter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("anonymous struct", func(t *testing.T) {
		p := struct {
			color string
		}{}
		assert.NoError(t, Set(&p, "color", "red"))
		assert.Equal(t, "red", p.color)
	})
	t.Run("anonymous *struct", func(t *testing.T) {
		p := &struct {
			color string
		}{}
		assert.NoError(t, Set(&p, "color", "brown"))
		assert.Equal(t, "brown", p.color)
	})
	t.Run("***struct", func(t *testing.T) {
		p := &struct {
			color string
		}{}
		p2 := &p
		p3 := &p2
		assert.NoError(t, Set(&p3, "color", "brown"))
		assert.Equal(t, "brown", p.color)
	})
	// todo
	//t.Run("var a interface{} = ***struct", func(t *testing.T) {
	//	p := &struct {
	//		color string
	//	}{}
	//	p2 := &p
	//	var p3 interface{} = &p2
	//	assert.NoError(t, Set(&p3, "color", "brown"))
	//	assert.Equal(t, "brown", p.color)
	//})
	t.Run("struct", func(t *testing.T) {
		p := person{}
		assert.NoError(t, Set(&p, "Name", "Jane"))
		assert.NoError(t, Set(&p, "age", 30))
		assert.Equal(t, person{Name: "Jane", age: 30}, p)
	})
	t.Run("*struct", func(t *testing.T) {
		p := &person{}
		assert.NoError(t, Set(&p, "Name", "Mary"))
		assert.NoError(t, Set(&p, "age", uint(33)))
		assert.Equal(t, &person{Name: "Mary", age: 33}, p)
	})
	t.Run("var a interface{} = &struct{}", func(t *testing.T) {
		var p interface{} = &person{}
		assert.NoError(t, Set(&p, "Name", "Mary Jane"))
		assert.NoError(t, Set(&p, "age", 45))
		assert.Equal(t, &person{Name: "Mary Jane", age: 45}, p)
	})
	t.Run("var a interface{} = struct{}", func(t *testing.T) {
		var p interface{} = person{}
		assert.NoError(t, Set(&p, "Name", "Jane"))
		assert.Equal(t, person{Name: "Jane"}, p)
	})
	t.Run("unexported type of field", func(t *testing.T) {
		p := person{}
		assert.NoError(t, Set(&p, "wallet", wallet{amount: 400}))
		assert.Equal(t, wallet{amount: 400}, p.wallet)
	})
	t.Run("zero values (val = nil)", func(t *testing.T) {
		t.Run("interface{}", func(t *testing.T) {
			p := struct {
				val interface{}
			}{val: 5}
			assert.NoError(t, Set(&p, "val", nil))
			assert.Nil(t, p.val)
		})
		t.Run("int", func(t *testing.T) {
			p := struct {
				val int
			}{val: 5}
			assert.NoError(t, Set(&p, "val", nil))
			assert.Equal(t, 0, p.val)
		})
		t.Run("struct", func(t *testing.T) {
			p := struct {
				val person
			}{
				val: person{
					Name: "Mary Jane",
					age:  35,
				},
			}
			assert.NoError(t, Set(&p, "val", nil))
			assert.Equal(t, person{}, p.val)
		})
		t.Run("unexported type of field", func(t *testing.T) {
			p := person{
				wallet: wallet{amount: 300},
			}
			assert.NoError(t, Set(&p, "wallet", nil))
			assert.Equal(t, wallet{}, p.wallet)
		})
	})
	t.Run("Given errors", func(t *testing.T) {
		t.Run("Field does not exist", func(t *testing.T) {
			p := person{}
			err := Set(&p, "FirstName", "Mary")
			assert.EqualError(t, err, "field `FirstName` does not exist")
		})
		t.Run("Invalid pointer dest", func(t *testing.T) {
			p := 5
			err := Set(&p, "FirstName", "Mary")
			assert.EqualError(t, err, "invalid parameter, setter.Set expects pointer to struct, given ptr.int")
		})
		t.Run("Invalid type of value", func(t *testing.T) {
			p := person{}
			err := Set(&p, "Name", struct{}{})
			assert.EqualError(t, err, "cannot cast `struct {}` to `string`")
		})
		t.Run("Invalid type of value (var p interface{} = person{})", func(t *testing.T) {
			var p interface{} = person{}
			err := Set(&p, "Name", struct{}{})
			assert.EqualError(t, err, "cannot cast `struct {}` to `string`")
		})
	})
}

func TestMustSet(t *testing.T) {
	t.Run("Given valid scenario", func(t *testing.T) {
		p := person{}
		MustSet(&p, "age", 75)
		assert.Equal(t, person{age: 75}, p)
	})
	t.Run("Given invalid scenario", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()
		MustSet(10, "foo", "bar")
	})
}

type person struct {
	Name   string
	age    uint8
	wallet wallet
}

type wallet struct {
	amount uint
}
