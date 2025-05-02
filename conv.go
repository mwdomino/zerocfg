package zerocfg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ToString returns a string representation of any value for use in configuration serialization and display.
//
// The conversion rules are as follows:
//   - If the value implements fmt.Stringer, its String() method is used.
//   - For slices or arrays whose elements implement fmt.Stringer, a JSON array of their string values is returned.
//   - For other slices, arrays, or maps, the value is marshaled to JSON.
//   - For all other types, fmt.Sprint is used.
//
// ToString is used internally by zerocfg for representing option values as strings, including when passing ToString to custom parsers and for rendering configuration output.
func ToString(v any) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}

	rv := reflect.ValueOf(v)
	rt := rv.Type()

	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		// check if element type implements fmt.Stringer
		stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
		if rt.Elem().Implements(stringerType) {
			n := rv.Len()
			strs := make([]string, n)
			for i := range strs {
				// safe to assert because we checked Implements
				strs[i] = rv.Index(i).Interface().(fmt.Stringer).String()
			}
			data, _ := json.Marshal(strs)
			return string(data)
		}

		// fallback: generic JSON marshal of the original slice/array
		fallthrough

	case reflect.Map:
		b, _ := json.Marshal(v)
		return string(b)

	default:
		// everything else: use fmt.Sprint
		return fmt.Sprint(v)
	}
}
