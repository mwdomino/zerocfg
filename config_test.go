package zerocfg

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testConfig() *config {
	return &config{
		make(map[string]*node),
		make(map[string]string),
		[]Parser{},
		false,
	}
}

func val[T any](v T, create func(T, *T) Value) Value {
	p := new(T)
	*p = v

	return create(v, p)
}

func Test_ConfigOk(t *testing.T) {
	const (
		name   = "name_key"
		alias  = "alias_key"
		desc   = "description"
		num    = 10
		prefix = "prefix"
	)

	tests := []struct {
		name   string
		setup  func()
		source map[string]any
		expect *config
	}{
		{
			name: "default",
			setup: func() {
				Int(name, num, desc)
				return
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
					},
				},
			},
		},
		{
			name: "override default",
			setup: func() {
				Int(name, 0, desc)
				return
			},
			source: map[string]any{
				name: num,
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						setSource:   mockType,
					},
				},
			},
		},
		{
			name: "alias",
			setup: func() {
				Int(name, 0, desc, Alias(alias))
				return
			},
			source: map[string]any{
				name: num,
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						Aliases:     []string{alias},
						setSource:   mockType,
					},
				},
				aliases: map[string]string{
					alias: name,
				},
			},
		},
		{
			name: "override alias",
			setup: func() {
				Int(name, 0, desc, Alias(alias))
				return
			},
			source: map[string]any{
				alias: num,
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						Aliases:     []string{alias},
						setSource:   mockType,
					},
				},
				aliases: map[string]string{
					alias: name,
				},
			},
		},
		{
			name: "group",
			setup: func() {
				g := NewGroup(prefix)
				Int(name, 0, desc, Group(g))
				return
			},
			source: map[string]any{
				prefix + "." + name: num,
			},
			expect: &config{
				vs: map[string]*node{
					prefix + "." + name: {
						Name:        prefix + "." + name,
						Description: desc,
						Value:       val(num, newIntValue),
						setSource:   mockType,
					},
				},
			},
		},
		{
			name: "option group",
			setup: func() {
				g := NewOptions(Secret())
				Int(name, 0, desc, Group(g))
				return
			},
			source: map[string]any{
				name: num,
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						setSource:   mockType,
						isSecret:    true,
					},
				},
			},
		},
		{
			name: "secret",
			setup: func() {
				Int(name, num, desc, Secret())
				return
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						isSecret:    true,
					},
				},
			},
		},
		{
			name: "required",
			setup: func() {
				Int(name, num, desc, Required())
				return
			},
			source: map[string]any{
				name: num,
			},
			expect: &config{
				vs: map[string]*node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						isRequired:  true,
						setSource:   mockType,
					},
				},
			},
		},
	}

	setConfig := func(expect *config) {
		expect.locked = true
		c.parsers = nil
		if expect.vs == nil {
			expect.vs = make(map[string]*node)
		}

		if expect.aliases == nil {
			expect.aliases = make(map[string]string)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = testConfig()
			tt.setup()

			err := Parse(newMock(tt.source))
			require.NoError(t, err)

			setConfig(tt.expect)
			require.EqualValues(t, tt.expect, c)
		})
	}
}

func Test_ConfigError(t *testing.T) {
	const (
		name    = "name_key"
		desc    = "description"
		unknown = "wrong_name"
	)

	var keyConflict = func(n1, n2 string, err error) error {
		return errorKeyConflict(
			&node{Name: n1},
			&node{Name: n2},
			err,
		)
	}

	tests := []struct {
		name    string
		setup   func()
		source  map[string]string
		err     error
		isPanic bool
	}{
		{
			name: "unknown option",
			setup: func() {
				Int(name, 0, desc)
				return
			},
			source: map[string]string{
				unknown: strconv.Itoa(0),
			},
			err: ErrNoSuchKey,
		},
		{
			name: "alias same as existing key",
			setup: func() {
				Int(name, 0, desc)
				Str(name+"1", "", desc, Alias(name))
				return
			},
			err:     keyConflict(name+"1", name, ErrCollidingAlias),
			isPanic: true,
		},
		{
			name: "key same as existing key",
			setup: func() {
				Int(name, 0, desc)
				Str(name, "", desc)
				return
			},
			err:     keyConflict(name, name, ErrDuplicateKey),
			isPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := func() error {
				c = defaultConfig()
				tt.setup()

				return c.applyParser(mockType, tt.source)
			}

			if tt.isPanic {
				require.PanicsWithError(t, tt.err.Error(), func() {
					_ = fn()
				})
			} else {
				require.ErrorIs(t, fn(), tt.err)
			}
		})
	}
}

func Test_ParseError(t *testing.T) {
	const (
		name = "name_key"
		desc = "description"
	)

	tests := []struct {
		name    string
		setup   func()
		source  map[string]any
		err     error
		isPanic bool
	}{
		{
			name: "missing required",
			setup: func() {
				Int(name, 0, desc, Required())
				return
			},
			source: map[string]any{},
			err:    ErrRequired,
		},
		{
			name: "double parse",
			setup: func() {
				Int(name, 0, desc, Required())
				_ = Parse()
				return
			},
			source: map[string]any{},
			err:    ErrDoubleParse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = testConfig()
			tt.setup()

			err := Parse(newMock(tt.source))
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func Test_RuntimeRegistration(t *testing.T) {
	c = testConfig()
	err := Parse(newMock(nil))
	require.NoError(t, err)

	const run = "runtime_key"
	require.PanicsWithError(t, fmt.Errorf("key=%q: %w", run, ErrRuntimeRegistration).Error(), func() {
		_ = Int(run, 0, "is not allowed")
	})
}

func Test_Render(t *testing.T) {
	c = testConfig()

	var names = []string{"a", "b", "c"}

	Uint(names[0], 1, "desc")
	Str(names[1], "", "desc")
	Int(names[2], 0, "desc")

	err := Parse(newMock(map[string]any{
		"a": 999,
	}))
	require.NoError(t, err)

	r := Show()
	for _, name := range names {
		require.True(t, strings.Contains(r, name))
	}

}

func Test_Priority(t *testing.T) {
	c = testConfig()

	const (
		dfault = "d"
		first  = "a"
		second = "b"
		third  = "c"
	)

	v0 := Int(dfault, 0, "")
	v1 := Int(first, 0, "")
	v2 := Int(second, 0, "")
	v3 := Int(third, 0, "")

	err := Parse(
		newMock(map[string]any{first: 1}),
		newMock(map[string]any{first: 2, second: 2}),
		newMock(map[string]any{first: 3, second: 3, third: 3}),
	)
	require.NoError(t, err)

	vs := []int{*v0, *v1, *v2, *v3}
	expected := []int{0, 1, 2, 3}
	require.Equal(t, expected, vs)
}

func Test_WrongType(t *testing.T) {
	c = testConfig()

	const (
		key   = "option_name"
		wrong = "wrong_value"
	)
	Int(key, 0, "")

	err := Parse(newMock(map[string]any{key: wrong}))
	require.Error(t, err)

	require.True(t, strings.Contains(err.Error(), key))
	require.True(t, strings.Contains(err.Error(), mockType))
	require.True(t, strings.Contains(err.Error(), wrong))
}
