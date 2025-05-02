package zerocfg

import (
	"encoding/json"
)

type mapValue map[string]any

func newMapValue(val map[string]any, p *map[string]any) Value {
	*p = val
	return (*mapValue)(p)
}

func (m *mapValue) Set(val string) error {
	for k := range *m {
		delete(*m, k)
	}

	return json.Unmarshal([]byte(val), m)
}

func (m *mapValue) Type() string {
	return "map"
}

func mapInternal(name string, value map[string]any, usage string, opts ...OptNode) *map[string]any {
	return Any(name, value, usage, newMapValue, opts...)
}

// Map registers a map[string]any configuration option and returns the map value.
//
// Arguments:
//   - name: unique option key (dot-separated for hierarchy)
//   - defVal: default map[string]any value
//   - desc: description for documentation and rendering
//   - opts: optional OptNode modifiers (e.g., Alias, Secret, Required)
//
// Returns:
//   - The registered map[string]any value, updated by configuration sources.
//
// Usage:
//
//	limits := zerocfg.Map("limits", map[string]any{"max": 10, "min": 1}, "map of limits")
func Map(name string, defVal map[string]any, desc string, opts ...OptNode) map[string]any {
	mptr := Any(name, defVal, desc, newMapValue, opts...)

	return *mptr
}
