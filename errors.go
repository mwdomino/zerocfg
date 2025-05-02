package zerocfg

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrNoSuchKey is returned when a configuration key is not registered as an option.
	ErrNoSuchKey = errors.New("no such key")

	// ErrCollidingAlias is returned when an alias collides with an existing key.
	ErrCollidingAlias = errors.New("colliding alias with key")

	// ErrDuplicateKey is returned when a duplicate configuration key is registered.
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrRequired is returned when required configuration fields are missing.
	ErrRequired = errors.New("missing required fields")

	// ErrRuntimeRegistration is returned when attempting to register options at runtime.
	ErrRuntimeRegistration = errors.New("misuse: runtime var registration is not allowed")

	// ErrDoubleParse is returned when Parse is called more than once.
	ErrDoubleParse = errors.New("misuse: Parse func should be called once")
)

// UnknownFieldError represents a mapping from configuration source names to unknown option keys encountered during parsing.
// It is returned by Parse when unknown values are found in configuration sources.
type UnknownFieldError map[string][]string

// IsUnknown checks if the provided error is an UnknownFieldError.
// If so, it returns the underlying map and true. Otherwise, it returns nil and false.
//
// Example usage:
//
//	err := zfg.Parse(...)
//	if u, ok := zfg.IsUnknown(err); ok {
//	    // u is map[source_name][]unknown_keys
//	}
func IsUnknown(err error) (map[string][]string, bool) {
	var v UnknownFieldError
	if !errors.As(err, &v) {
		return nil, false
	}
	return v, true
}

func (e UnknownFieldError) Error() string {
	data, _ := json.Marshal(e)

	return fmt.Sprintf("unknown fields: %s", string(data))
}

func (e *UnknownFieldError) add(source string, unknown map[string]string) {
	if len(unknown) == 0 {
		return
	}

	s := make([]string, 0, len(unknown))
	for k := range unknown {
		s = append(s, k)
	}

	(*e)[source] = s
}
