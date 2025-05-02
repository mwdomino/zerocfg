package zerocfg

import (
	"fmt"
	"strings"
)

type Parser interface {
	Type() string
	Parse(awaited map[string]bool, conv func(any) string) (found, unknown map[string]string, err error)
}

func Parse(ps ...Parser) error {
	if c.locked {
		return ErrDoubleParse
	}
	c.locked = true
	c.parsers = append(c.parsers, ps...)
	awaited := c.awaited()

	uErr := make(UnknownFieldError)
	for _, p := range c.parsers {
		found, unknown, err := p.Parse(awaited, ToString)
		if err != nil {
			return fmt.Errorf("parse %q: %w", p.Type(), err)
		}

		err = c.applyParser(found)
		if err != nil {
			return fmt.Errorf("apply %q: %w", p.Type(), err)
		}

		uErr.add(p.Type(), unknown)
	}

	if len(uErr) != 0 {
		return uErr
	}

	var required []string
	for _, v := range c.vs {
		if v.isRequired && !v.fromSource {
			required = append(required, v.Name)
		}
	}
	if len(required) != 0 {
		return fmt.Errorf("%w: %s", ErrRequired, strings.Join(required, ", "))
	}

	return nil
}

func (c *config) applyParser(vs map[string]string) error {
	for k, v := range vs {
		err := c.set(k, v)
		if err != nil {
			return fmt.Errorf("set key=%q: %w", k, err)
		}
	}

	return nil
}
