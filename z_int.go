package zerocfg

import (
	"encoding/json"
	"strconv"
)

type intValue int

func newIntValue(val int, p *int) Value {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Type() string {
	return "int"
}

// Int registers an int configuration option and returns a pointer to its value.
//
// Arguments:
//   - name: unique option key (dot-separated for hierarchy)
//   - defVal: default int value
//   - desc: description for documentation and rendering
//   - opts: optional OptNode modifiers (e.g., Alias, Secret, Required)
//
// Returns:
//   - Pointer to the registered int value, updated by configuration sources.
//
// Usage:
//
//	port := zerocfg.Int("db.port", 5432, "database port")
func Int(name string, defVal int, desc string, opts ...OptNode) *int {
	return Any(name, defVal, desc, newIntValue, opts...)
}

type int32Value int32

func newInt32Value(val int32, p *int32) Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) Type() string {
	return "int32"
}

// Int32 registers an int32 configuration option and returns a pointer to its value.
//
// Arguments:
//   - name: unique option key (dot-separated for hierarchy)
//   - defVal: default int32 value
//   - desc: description for documentation and rendering
//   - opts: optional OptNode modifiers (e.g., Alias, Secret, Required)
//
// Returns:
//   - Pointer to the registered int32 value, updated by configuration sources.
//
// Usage:
//
//	code := zerocfg.Int32("status.code", 200, "status code")
func Int32(name string, defVal int32, desc string, opts ...OptNode) *int32 {
	return Any(name, defVal, desc, newInt32Value, opts...)
}

type int64Value int64

func newInt64Value(val int64, p *int64) Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Type() string {
	return "int64"
}

// Int64 registers an int64 configuration option and returns a pointer to its value.
//
// Arguments:
//   - name: unique option key (dot-separated for hierarchy)
//   - defVal: default int64 value
//   - desc: description for documentation and rendering
//   - opts: optional OptNode modifiers (e.g., Alias, Secret, Required)
//
// Returns:
//   - Pointer to the registered int64 value, updated by configuration sources.
//
// Usage:
//
//	big := zerocfg.Int64("big.value", 1234567890, "big int value")
func Int64(name string, defVal int64, desc string, opts ...OptNode) *int64 {
	return Any(name, defVal, desc, newInt64Value, opts...)
}

type intSliceValue []int

func newIntSlice(val []int, p *[]int) Value {
	*p = val
	return (*intSliceValue)(p)
}

func (s *intSliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *intSliceValue) Type() string {
	return "ints"
}

// Ints registers a slice of int configuration options and returns a pointer to its value.
//
// Arguments:
//   - name: unique option key (dot-separated for hierarchy)
//   - defVal: default slice of int values
//   - desc: description for documentation and rendering
//   - opts: optional OptNode modifiers (e.g., Alias, Secret, Required)
//
// Returns:
//   - Pointer to the registered slice of int values, updated by configuration sources.
//
// Usage:
//
//	ids := zerocfg.Ints("user.ids", []int{1, 2, 3}, "user IDs")
func Ints(name string, defVal []int, desc string, opts ...OptNode) *[]int {
	return Any(name, defVal, desc, newIntSlice, opts...)
}
