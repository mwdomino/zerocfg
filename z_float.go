package zerocfg

import (
	"encoding/json"
	"strconv"
)

type float64Value float64

func newFloat64(val float64, p *float64) Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(val string) error {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	*f = float64Value(v)
	return nil
}

func (f *float64Value) Type() string {
	return "float64"
}

func Float64(name string, value float64, usage string, opts ...OptNode) *float64 {
	return Any(name, value, usage, newFloat64, opts...)
}

type float64SliceValue []float64

func newFloat64Slice(val []float64, p *[]float64) Value {
	*p = val
	return (*float64SliceValue)(p)
}

func (s *float64SliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *float64SliceValue) Type() string {
	return "floats64"
}

func Floats64(name string, value []float64, usage string, opts ...OptNode) *[]float64 {
	return Any(name, value, usage, newFloat64Slice, opts...)
}

type float32Value float32

func newFloat32(val float32, p *float32) Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(val string) error {
	v, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return err
	}
	*f = float32Value(float32(v))
	return nil
}

func (f *float32Value) Type() string {
	return "float32"
}

func Float32(name string, value float32, usage string, opts ...OptNode) *float32 {
	return Any(name, value, usage, newFloat32, opts...)
}

type float32SliceValue []float32

func newFloat32Slice(val []float32, p *[]float32) Value {
	*p = val
	return (*float32SliceValue)(p)
}

func (s *float32SliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *float32SliceValue) Type() string {
	return "floats32"
}

func Floats32(name string, value []float32, usage string, opts ...OptNode) *[]float32 {
	return Any(name, value, usage, newFloat32Slice, opts...)
}
