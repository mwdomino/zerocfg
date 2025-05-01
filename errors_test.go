package zerocfg

import (
	"io"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UnknownOk(t *testing.T) {
	c = testConfig()
	Str("known", "", "")

	p := newMock(map[string]any{
		"v1":       1,
		"v2":       "",
		"v3.a.b.c": false,
		"known":    "field",
	})

	err := Parse(p)
	u, ok := IsUnknown(err)
	require.True(t, ok)

	expected := UnknownFieldError{
		mockType: []string{"v1", "v2", "v3.a.b.c"},
	}

	sort.Strings(expected[mockType])
	sort.Strings(u[mockType])
	require.EqualValues(t, expected, u)
	require.Equal(t, expected.Error(), err.Error())
}

func Test_NotUnknown(t *testing.T) {
	_, ok := IsUnknown(io.ErrClosedPipe)
	require.False(t, ok)
}
