package caller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustCall(t *testing.T) {
	t.Run("Given method", func(t *testing.T) {
		p := person{}
		assert.Equal(t, "", p.name)
		MustCall(p.setName, "Jane")
		assert.Equal(t, "Jane", p.name)
	})

	t.Run("Given invalid functions", func(t *testing.T) {
		scenarios := []struct {
			fn interface{}
		}{
			{fn: 5},
			{fn: false},
			{fn: (*error)(nil)},
			{fn: struct{}{}},
		}
		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				defer func() {
					v := recover()
					assert.NotNil(t, v)
					assert.Regexp(t, "^expects func, .* given$", v)
				}()
				MustCall(s.fn)
			})
		}
	})

	t.Run("Given invalid argument", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "arg0: cannot cast `struct {}` to `[]int`", fmt.Sprintf("%s", recover()))
		}()

		MustCall(func([]int) {}, struct{}{})
	})

	t.Run("Given too many arguments", func(t *testing.T) {
		scenarios := []struct {
			fn   interface{}
			args []interface{}
		}{
			{
				fn:   strings.Join,
				args: []interface{}{"1", "2", "3"},
			},
			{
				fn:   func() {},
				args: []interface{}{1},
			},
		}
		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				defer func() {
					v := recover()
					assert.NotNil(t, v)
					assert.Equal(t, "call with too many input arguments", v)
				}()
				MustCall(s.fn, s.args...)
			})
		}
	})

	t.Run("Given too few arguments", func(t *testing.T) {
		scenarios := []struct {
			fn   interface{}
			args []interface{}
		}{
			{
				fn:   strings.Join,
				args: []interface{}{},
			},
			{
				fn:   func(a int) {},
				args: []interface{}{},
			},
		}
		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				defer func() {
					v := recover()
					assert.NotNil(t, v)
					assert.Equal(t, "call with too few input arguments", v)
				}()
				MustCall(s.fn, s.args...)
			})
		}
	})

	t.Run("Given scenarios", func(t *testing.T) {
		scenarios := []struct {
			fn       interface{}
			args     []interface{}
			expected []interface{}
		}{
			{
				fn: func(a, b int) int {
					return a + b
				},
				args:     []interface{}{uint(1), uint(2)},
				expected: []interface{}{int(3)},
			},
			{
				fn: func(a, b uint) uint {
					return a + b
				},
				args:     []interface{}{int(7), int(3)},
				expected: []interface{}{uint(10)},
			},
			{
				fn: func(vals ...uint) (result uint) {
					for _, v := range vals {
						result += v
					}
					return
				},
				args:     []interface{}{int(1), int(2), int(3)},
				expected: []interface{}{uint(6)},
			},
		}
		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				assert.Equal(
					t,
					s.expected,
					MustCall(s.fn, s.args...),
				)
			})
		}
	})

	t.Run("Convert parameters", func(t *testing.T) {
		scenarios := map[string]struct {
			fn     interface{}
			input  interface{}
			output interface{}
			error  string
		}{
			"[]interface{} to []type": {
				fn: func(v []int) []int {
					return v
				},
				input:  []interface{}{1, 2, 3},
				output: []int{1, 2, 3},
			},
			"[]struct{}{} to []type": {
				fn:    func([]int) {},
				input: []struct{}{},
				error: "arg0: cannot cast `[]struct {}` to `[]int`",
			},
			"nil to interface{}": {
				fn: func(v interface{}) interface{} {
					return v
				},
				input:  nil,
				output: (interface{})(nil),
			},
		}

		for n, s := range scenarios {
			t.Run(n, func(t *testing.T) {
				defer func() {
					if s.error == "" {
						assert.Nil(t, recover())
						return
					}

					assert.Equal(
						t,
						s.error,
						fmt.Sprintf("%s", recover()),
					)
				}()

				r := MustCall(s.fn, s.input)
				assert.Len(t, r, 1)
				assert.Equal(t, r[0], s.output)
			})
		}
	})
}

func TestCall(t *testing.T) {
	t.Run("Given scenario", func(t *testing.T) {
		r, err := Call(fmt.Sprintf, "%s %s", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(
			t,
			[]interface{}{"hello world"},
			r,
		)
	})
	t.Run("Given error", func(t *testing.T) {
		r, err := Call(fmt.Sprintf)
		assert.Error(t, err)
		assert.Nil(t, r)
	})
}

func TestCallProvider(t *testing.T) {
	t.Run("Given scenarios", func(t *testing.T) {
		scenarios := []struct {
			provider interface{}
			params   []interface{}
			expected interface{}
		}{
			{
				provider: func() interface{} {
					return nil
				},
				expected: nil,
			},
			{
				provider: func(vals ...int) (int, error) {
					result := 0
					for _, v := range vals {
						result += v
					}

					return result, nil
				},
				params:   []interface{}{10, 100, 200},
				expected: 310,
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				r, err := CallProvider(s.provider, s.params...)
				assert.NoError(t, err)
				assert.Equal(t, s.expected, r)

				t.Run("MustCallProvider", func(t *testing.T) {
					defer func() {
						assert.Nil(t, recover())
					}()
					assert.Equal(t, s.expected, MustCallProvider(s.provider, s.params...))
				})
			})
		}
	})

	t.Run("Given errors", func(t *testing.T) {
		scenarios := []struct {
			provider interface{}
			params   []interface{}
			err      string
		}{
			{
				provider: func() {},
				err:      "provider must return 1 or 2 values, given function returns 0 values",
			},
			{
				provider: func() (interface{}, interface{}, interface{}) {
					return nil, nil, nil
				},
				err: "provider must return 1 or 2 values, given function returns 3 values",
			},
			{
				provider: func() (interface{}, interface{}) {
					return nil, nil
				},
				err: "second value returned by provider must implements error interface",
			},
			{
				provider: func() (interface{}, error) {
					return nil, fmt.Errorf("test error")
				},
				err: "test error",
			},
			{
				provider: func() interface{} {
					return nil
				},
				params: []interface{}{1, 2, 3},
				err:    "call with too many input arguments",
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				r, err := CallProvider(s.provider, s.params...)
				assert.Nil(t, r)
				assert.EqualError(t, err, s.err)

				t.Run("MustCallProvider", func(t *testing.T) {
					defer func() {
						assert.Equal(
							t,
							s.err,
							fmt.Sprintf("%s", recover()),
						)
					}()
					MustCallProvider(s.provider, s.params...)
				})
			})
		}
	})

	t.Run("Given invalid provider", func(t *testing.T) {
		_, err := CallProvider(5)
		assert.EqualError(t, err, "provider must be kind of func, int given")

		t.Run("MustCallProvider", func(t *testing.T) {
			defer func() {
				assert.Equal(
					t,
					"provider must be kind of func, int given",
					fmt.Sprintf("%s", recover()),
				)
			}()
			MustCallProvider(5)
		})
	})
}

func TestWrapMustCallProvider(t *testing.T) {
	t.Run("Given scenario", func(t *testing.T) {
		defer assert.Nil(t, recover())
		expected := person{name: "Jane"}
		r := WrapMustCallProvider("custom error", func() person {
			return expected
		})
		assert.Equal(t, expected, r)
	})
	t.Run("Given error", func(t *testing.T) {
		defer func() {
			assert.Equal(
				t,
				"custom error: provider must be kind of func, int given",
				fmt.Sprintf("%s", recover()),
			)
		}()
		WrapMustCallProvider("custom error", 5)
	})
}

func TestCallWitherByName(t *testing.T) {
	t.Run("Given scenarios", func(t *testing.T) {
		var emptyPerson interface{} = person{}

		scenarios := []struct {
			object interface{}
			wither string
			params []interface{}
			output interface{}
		}{
			{
				object: make(ints, 0),
				wither: "Append",
				params: []interface{}{5},
				output: ints{5},
			},
			{
				object: person{name: "Mary"},
				wither: "WithName",
				params: []interface{}{"Jane"},
				output: person{name: "Jane"},
			},
			{
				object: &person{name: "Mary"},
				wither: "WithName",
				params: []interface{}{"Jane"},
				output: person{name: "Jane"},
			},
			{
				object: emptyPerson,
				wither: "WithName",
				params: []interface{}{"Kaladin"},
				output: person{name: "Kaladin"},
			},
			{
				object: &emptyPerson,
				wither: "WithName",
				params: []interface{}{"Shallan"},
				output: person{name: "Shallan"},
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				result, err := CallWitherByName(s.object, s.wither, s.params...)
				assert.NoError(t, err)
				assert.Equal(t, s.output, result)
			})
		}
	})

	t.Run("Given errors", func(t *testing.T) {
		scenarios := []struct {
			object interface{}
			wither string
			params []interface{}
			error  string
		}{
			{
				object: person{},
				wither: "withName",
				params: nil,
				error:  "invalid wither `caller.person`.`withName`",
			},
			{
				object: person{},
				wither: "Clone",
				params: nil,
				error:  "wither must return 1 value, given function returns 2 values",
			},
		}

		for i, s := range scenarios {
			t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
				o, err := CallWitherByName(s.object, s.wither, s.params...)
				assert.Nil(t, o)
				assert.EqualError(t, err, s.error)
			})
		}
	})
}

type ints []int

func (i ints) Append(v int) ints {
	return append(i, v)
}

type person struct {
	name string
}

func (p person) Clone() (person, error) {
	return p, nil
}

func (p person) WithName(n string) person {
	return person{name: n}
}

func (p person) withName(n string) person {
	return person{name: n}
}

func (p *person) setName(n string) {
	p.name = n
}
