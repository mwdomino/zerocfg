package zerocfg

import (
	"fmt"
	"strings"
)

// Parser defines a configuration source for zerocfg.
//
// Custom sources must implement this interface to provide configuration values.
//
// Methods:
//   - Type() string: returns the parser's type name (e.g., "env", "yaml").
//   - Parse(awaited, conv):
//   - awaited: map of option names and aliases to expect (true = option, false = alias)
//   - conv: function to convert values to string (usually zerocfg.ToString)
//
// Returns:
//   - found: map of recognized option names to string values
//   - unknown: map of unrecognized names to string values
type Parser interface {
	Type() string
	Parse(awaited map[string]bool, conv func(any) string) (found, unknown map[string]string, err error)
}

// Parse loads configuration from the provided sources in priority order.
//
// Usage:
//
//	err := zerocfg.Parse(env.New(), yaml.New(path))
//
// Priority:
//  1. Command-line flags (always highest)
//  2. Parsers in the order provided (first = higher priority)
//  3. Default values (lowest)
//
// Behavior:
//   - Applies each parser in order, setting values for registered options only.
//   - Returns an error if unknown options are found (unless ignored by IsUnknown).
//
// Error Handling:
//   - UnknownFieldError: for unknown keys (see IsUnknown)
//   - ErrRequired: for missing required options
//   - ErrDoubleParse: if called multiple times
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
