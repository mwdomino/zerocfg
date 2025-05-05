package zerocfg

const mockType = "mock"

type mockParser struct {
	values map[string]any
}

func newMock(v map[string]any) *mockParser {
	return &mockParser{values: v}
}

func (m mockParser) Type() string {
	return mockType
}

func (m mockParser) Provide(awaited map[string]bool, conv func(any) string) (f, u map[string]string, _ error) {
	f, u = map[string]string{}, map[string]string{}
	for k, v := range m.values {
		if _, ok := awaited[k]; ok {
			f[k] = conv(v)
		} else {
			u[k] = conv(v)
		}
	}

	return
}
