package zerocfg

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNoSuchKey           = errors.New("no such key")
	ErrCollidingAlias      = errors.New("colliding alias with key")
	ErrDuplicateKey        = errors.New("duplicate key")
	ErrRequired            = errors.New("missing required fields")
	ErrRuntimeRegistration = errors.New("misuse: runtime var registration is not allowed")
	ErrDoubleParse         = errors.New("misuse: Parse func should be called once")
)

type UnknownFieldError map[string][]string

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
