package zerocfg

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const name = "sample.name"

type anyFn[T any] func(string, T, string, ...OptNode) *T

func regSource[T any](fn anyFn[T], v T) (reg func() any, value any, source map[string]string) {
	var def T
	reg = func() any {
		return fn(name, def, "desc")
	}

	source = map[string]string{name: ToString(v)}

	return reg, v, source
}

func Test_ValueOk(t *testing.T) {
	tests := []struct {
		varType string
		init    func() (ptr func() any, val any, src map[string]string)
	}{
		{
			varType: "int",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int, 5)
			},
		},
		{
			varType: "uint",
			init: func() (func() any, any, map[string]string) {
				return regSource(Uint, uint(42))
			},
		},
		{
			varType: "int32",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int32, int32(123))
			},
		},
		{
			varType: "uint32",
			init: func() (func() any, any, map[string]string) {
				return regSource(Uint32, uint32(456))
			},
		},
		{
			varType: "int64",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int64, int64(789))
			},
		},
		{
			varType: "uint64",
			init: func() (func() any, any, map[string]string) {
				return regSource(Uint64, uint64(1011))
			},
		},
		{
			varType: "bool true",
			init: func() (func() any, any, map[string]string) {
				return regSource(Bool, true)
			},
		},
		{
			varType: "bool false",
			init: func() (func() any, any, map[string]string) {
				return regSource(Bool, false)
			},
		},
		{
			varType: "bools",
			init: func() (func() any, any, map[string]string) {
				return regSource(Bools, []bool{true, false, true})
			},
		},
		{
			varType: "ints",
			init: func() (func() any, any, map[string]string) {
				return regSource(Ints, []int{1, 2, 3})
			},
		},
		{
			varType: "string",
			init: func() (func() any, any, map[string]string) {
				return regSource(Str, "value")
			},
		},
		{
			varType: "strings",
			init: func() (func() any, any, map[string]string) {
				return regSource(Strs, []string{"a", "b", "c"})
			},
		},
		{
			varType: "float32",
			init: func() (func() any, any, map[string]string) {
				return regSource(Float32, float32(3.14))
			},
		},
		{
			varType: "float32s",
			init: func() (func() any, any, map[string]string) {
				return regSource(Floats32, []float32{1.1, 2.2, 3.3})
			},
		},
		{
			varType: "float64",
			init: func() (func() any, any, map[string]string) {
				return regSource(Float64, 3.14159265359)
			},
		},
		{
			varType: "float64s",
			init: func() (func() any, any, map[string]string) {
				return regSource(Floats64, []float64{1.1, 2.2, 3.3})
			},
		},
		{
			varType: "duration",
			init: func() (func() any, any, map[string]string) {
				return regSource(Dur, 5*time.Second)
			},
		},
		{
			varType: "durations",
			init: func() (func() any, any, map[string]string) {
				return regSource(Durs, []time.Duration{time.Second, 2 * time.Minute, 3 * time.Hour})
			},
		},
		{
			varType: "ip",
			init: func() (func() any, any, map[string]string) {
				return regSource(ipInternal, net.ParseIP("192.168.1.1"))
			},
		},
		{
			varType: "map",
			init: func() (func() any, any, map[string]string) {
				return regSource(mapInternal, map[string]any{"float": 1., "str": "val"})
			},
		},
	}

	dereference := func(v any) (any, error) {
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("not a pointer")
		}
		return val.Elem().Interface(), nil
	}

	for _, tt := range tests {
		t.Run(tt.varType, func(t *testing.T) {
			c = defaultConfig()

			reg, expected, source := tt.init()
			ptr := reg()

			err := c.applyParser(source)
			require.NoError(t, err)

			actual, err := dereference(ptr)
			require.NoError(t, err)

			require.EqualValues(t, expected, actual)

			// check string representation <-> sources map[string]string
			node, ok := c.vs[name]
			require.True(t, ok)

			newValue := node.Value.String()
			source[name] = newValue

			err = c.applyParser(source)
			require.NoError(t, err)

			updatedNode, ok := c.vs[name]
			require.True(t, ok)

			require.Equal(t, newValue, updatedNode.Value.String())

			// check type name
			awaitedType := strings.Split(tt.varType, " ")[0]
			require.Equal(t, awaitedType, updatedNode.Value.Type())
		})
	}
}
