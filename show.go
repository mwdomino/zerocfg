package zerocfg

import (
	"bytes"
	"fmt"
	"sort"
)

// Show returns a formatted string representation of all registered configuration options and their current values.
func Show() string {
	vs := make([]*node, 0, len(c.vs))
	for _, n := range c.vs {
		vs = append(vs, n)
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Name < vs[j].Name
	})

	return render(vs)
}

func render(vs []*node) string {
	var maxName, maxVal int
	for _, v := range vs {
		if l := len(v.Name); l > maxName {
			maxName = l
		}

		val := ToString(v.Value)
		if v.isSecret {
			val = "<secret>"
		}

		if l := len(val); l > maxVal {
			maxVal = l
		}
	}

	var b bytes.Buffer
	for _, v := range vs {
		val := ToString(v.Value)
		if v.isSecret {
			val = "<secret>"
		}

		line := fmt.Sprintf(
			"%-*s = %-*s (%s)\n",
			maxName, v.Name,
			maxVal, val,
			v.Description,
		)

		b.WriteString(line)
	}

	return b.String()
}
