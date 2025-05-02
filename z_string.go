package zerocfg

import "encoding/json"

type stringValue string

func newStringValue(val string, p *string) Value {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Type() string {
	return "string"
}

func Str(name string, defVal string, desc string, opts ...OptNode) *string {
	return Any(name, defVal, desc, newStringValue, opts...)
}

type stringSliceValue []string

func newStringSlice(val []string, p *[]string) Value {
	*p = val
	return (*stringSliceValue)(p)
}

func (s *stringSliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *stringSliceValue) Type() string {
	return "strings"
}

func Strs(name string, defVal []string, desc string, opts ...OptNode) *[]string {
	return Any(name, defVal, desc, newStringSlice, opts...)
}
