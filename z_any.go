package zerocfg

import "fmt"

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
