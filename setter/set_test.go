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
	t.Run("var a interface{}", func(t *testing.T) {
		t.Run("&struct{}", func(t *testing.T) {
			const color = "red"
			p := struct {
				color string
			}{}
			var obj interface{} = &p
			assert.Equal(t, "", p.color)
			assert.NoError(t, Set(obj, "color", color))
			assert.Equal(t, color, p.color)
		})
		t.Run("&&struct{}", func(t *testing.T) {
			const color = "blue"
			p := struct {
				color string
			}{}
			p2 := &p
			var obj interface{} = &p2
			assert.Equal(t, "", p.color)
			assert.NoError(t, Set(obj, "color", color))
			assert.Equal(t, color, p.color)
		})
		t.Run("&&&struct{}", func(t *testing.T) {
			const color = "yellow"
			p := struct {
				color string
			}{}
			p2 := &p
			p3 := &p2
			var obj interface{} = &p3
			assert.Equal(t, "", p.color)
			assert.NoError(t, Set(obj, "color", color))
			assert.Equal(t, color, p.color)
		})
		t.Run("&&&&struct{}", func(t *testing.T) {
			const color = "green"
			p := struct {
				color string
			}{}
			p2 := &p
			p3 := &p2
			p4 := &p3
			var obj interface{} = &p4
			assert.Equal(t, "", p.color)
			assert.NoError(t, Set(obj, "color", color))
			assert.Equal(t, color, p.color)
		})
	})
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
	t.Run("convert []interface{} to []type", func(t *testing.T) {
		s := storage{}
		assert.NoError(
			t,
			Set(&s, "wallets", []interface{}{wallet{100}, wallet{200}}),
		)
		assert.Equal(
			t,
			[]wallet{{100}, {200}},
			s.wallets,
		)
	})
	t.Run("Given errors", func(t *testing.T) {
		t.Run("Field does not exist", func(t *testing.T) {
			p := person{}
			err := Set(&p, "FirstName", "Mary")
			assert.EqualError(t, err, "set `*setter.person`.`FirstName`: field `FirstName` does not exist")
		})
		t.Run("Invalid pointer dest", func(t *testing.T) {
			p := 5
			err := Set(&p, "FirstName", "Mary")
			assert.EqualError(t, err, "invalid parameter, setter.Set expects pointer to struct, given ptr.int")
		})
		t.Run("Invalid type of value", func(t *testing.T) {
			p := person{}
			err := Set(&p, "Name", struct{}{})
			assert.EqualError(t, err, "set `*setter.person`.`Name`: cannot cast `struct {}` to `string`")
		})
		t.Run("Invalid type of value (var p interface{} = person{})", func(t *testing.T) {
			var p interface{} = person{}
			err := Set(&p, "Name", struct{}{})
			assert.EqualError(t, err, "set `*interface {}`.`Name`: cannot cast `struct {}` to `string`")
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

type storage struct {
	wallets []wallet
}
