package zerocfg

import "strconv"

type uintValue uint

func newUintValue(val uint, p *uint) Value {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Type() string {
	return "uint"
}

// Uint registers a uint configuration option and returns a pointer to its value.
//
// Usage:
//
//	port := zerocfg.Uint("db.port", 5678, "database port")
func Uint(name string, defVal uint, desc string, opts ...OptNode) *uint {
	return Any(name, defVal, desc, newUintValue, opts...)
}

type uint32Value uint32

func newUint32Value(val uint32, p *uint32) Value {
	*p = val
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
}

func (i *uint32Value) Type() string {
	return "uint32"
}

// Uint32 registers a uint32 configuration option and returns a pointer to its value.
//
// Usage:
//
//	code := zerocfg.Uint32("status.code", 200, "status code")
func Uint32(name string, defVal uint32, desc string, opts ...OptNode) *uint32 {
	return Any(name, defVal, desc, newUint32Value, opts...)
}

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Type() string {
	return "uint64"
}

// Uint64 registers a uint64 configuration option and returns a pointer to its value.
//
// Usage:
//
//	big := zerocfg.Uint64("big.value", 1234567890, "big uint value")
func Uint64(name string, defVal uint64, desc string, opts ...OptNode) *uint64 {
	return Any(name, defVal, desc, newUint64Value, opts...)
}
