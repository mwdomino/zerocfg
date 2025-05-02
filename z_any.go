package zerocfg

import "fmt"

// Any registers a custom configuration option type and returns a pointer to its value.
//
// Usage Example:
//
//	// Custom type
//	type MyType struct { V string }
//
//	func newValue(val MyType, p *MyType) Value {
//	    *p = val
//	    return (*MyType)(p)
//	}
//
//	func (m *MyType) Set(s string) error { m.V = s; return nil }
//	func (m *MyType) Type() string      { return "custom" }
//
//	// User-friendly registration function
//	func Custom(name string, defVal MyType, desc string, opts ...zfg.OptNode) *MyType {
//	    return zfg.Any(name, defVal, desc, newValue, opts...)
//	}
//
//	// Register custom option
//	myOpt := Custom("custom.opt", MyType{"default"}, "custom option")
//
// Arguments:
//   - name: unique option key (dot-separated for hierarchy)
//   - defVal: default value of type T
//   - desc: description for documentation and rendering
//   - create: function to create a Value implementation for T (see example above)
//   - opts: optional OptNode modifiers (e.g., Alias, Secret, Required)
//
// Behavior:
//   - Registers the option at import time; panics if called after Parse.
//   - Returns a pointer to the registered value, which is updated by configuration sources.
func Any[T any](name string, defVal T, desc string, create func(T, *T) Value, opts ...OptNode) *T {
	if c.locked {
		err := fmt.Errorf("key=%q: %w", name, ErrRuntimeRegistration)
		panic(err)
	}

	p := new(T)
	*p = defVal
	c.add(name, create(defVal, p), desc, opts...)

	return p
}
