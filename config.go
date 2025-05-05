package zerocfg

import (
	"fmt"

	"github.com/chaindead/zerocfg/flag"
)

type config struct {
	vs      map[string]*node
	aliases map[string]string

	parsers []Provider
	locked  bool
}

func defaultConfig() *config {
	return &config{
		make(map[string]*node),
		make(map[string]string),
		[]Provider{flag.New()},
		false,
	}
}

var c = defaultConfig()

func (c *config) add(key string, v Value, usage string, opts ...OptNode) {
	n := &node{
		Name:        key,
		Description: usage,
		Value:       v,
		caller:      findCaller(),
	}

	for _, opt := range opts {
		opt(n)
	}

	if existing, ok := c.vs[n.Name]; ok {
		err := errorKeyConflict(n, existing, ErrDuplicateKey)
		panic(err)
	}

	c.vs[n.Name] = n
	for _, alias := range n.Aliases {
		if existing, ok := c.vs[alias]; ok {
			err := errorKeyConflict(n, existing, ErrCollidingAlias)
			panic(err)
		}

		c.aliases[alias] = n.Name
	}
}

func errorKeyConflict(new *node, existing *node, err error) error {
	return fmt.Errorf("key %q confilicts with %q: %w", new.pathName(), existing.pathName(), err)
}

func (c *config) set(source, key string, v string) error {
	trueKey, ok := c.aliases[key]
	if ok {
		key = trueKey
	}

	n, ok := c.vs[key]
	if !ok {
		return ErrNoSuchKey
	}

	if n.setSource != "" {
		return nil
	}

	n.setSource = source
	return c.vs[key].Value.Set(v)
}

func (c *config) awaited() map[string]bool {
	a := make(map[string]bool)

	for k := range c.vs {
		a[k] = true
	}

	for k := range c.aliases {
		a[k] = false
	}

	return a
}
