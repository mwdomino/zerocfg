package env

import (
	"os"
	"regexp"
	"strings"
)

var cleanRe = regexp.MustCompile(`[^A-Za-z0-9.]+`)

type Opt func(*Parser)

// WithPrefix returns an Opt that sets the prefix for environment variable names in the Parser.
func WithPrefix(prefix string) Opt {
	return func(p *Parser) {
		p.prefix = prefix
	}
}

// Parser parses environment variables for configuration.
type Parser struct {
	// Prefix to prepend to all environment variable names.
	prefix string
}

// New creates a new Parser with the provided options.
func New(opts ...Opt) *Parser {
	p := &Parser{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Type returns the type name of the parser.
func (p Parser) Type() string {
	return "env"
}

// Parse reads environment variables matching the awaited keys and returns found values.
func (p Parser) Parse(awaited map[string]bool, _ func(any) string) (found, unknown map[string]string, err error) {
	keys := make(map[string]string, len(awaited))
	for k := range awaited {
		keys[p.prefix+k] = toENV(k)
	}

	found = make(map[string]string)
	for original, formatted := range keys {
		v, ok := os.LookupEnv(formatted)
		if !ok {
			continue
		}

		found[original] = v
	}

	return found, unknown, nil
}

// toENV transforms the input string into an uppercase, underscore-separated
// environment variable name by:
// 1. Removing all characters except letters, digits, and dots.
// 2. Converting to uppercase.
// 3. Replacing dots with underscores.
func toENV(s string) string {
	cleaned := cleanRe.ReplaceAllString(s, "")
	upper := strings.ToUpper(cleaned)
	envName := strings.ReplaceAll(upper, ".", "_")

	return envName
}
