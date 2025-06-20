package zerocfg

import (
	"bytes"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
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

// AsYaml returns a YAML representation of the complete configuration.
func AsYaml() (string, error) {
	vs := make([]*node, 0, len(c.vs))
	for _, n := range c.vs {
		vs = append(vs, n)
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Name < vs[j].Name
	})

	config := buildConfigMap(vs)
	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

func buildConfigMap(vs []*node) map[string]any {
	config := make(map[string]any)

	for _, v := range vs {
		setNestedValue(config, v.Name, yamlConfigValue(v))
	}

	return config
}

func setNestedValue(config map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	current := config

	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
		} else {
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]any)
			}
			current = current[part].(map[string]any)
		}
	}
}

func yamlConfigValue(n *node) any {
	if n.isSecret {
		return "<secret>"
	}

	return ToString(n.Value)
}

func render(vs []*node) string {
	doc := &yaml.Node{Kind: yaml.DocumentNode}
	root := &yaml.Node{Kind: yaml.MappingNode}
	doc.Content = []*yaml.Node{root}

	for _, v := range vs {
		addNode(root, v.Name, yamlValue(v), yamlDescription(v))
	}

	var buf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2)
	_ = yamlEncoder.Encode(doc)
	_ = yamlEncoder.Close()
	return buf.String()
}

func addNode(root *yaml.Node, path string, value any, comment string) {
	parts := strings.Split(path, ".")

	cur := root
	for i, part := range parts {
		var keyNode *yaml.Node
		var valueNode *yaml.Node
		found := false
		for j := 0; j+1 < len(cur.Content); j += 2 {
			if cur.Content[j].Value == part {
				keyNode = cur.Content[j]
				valueNode = cur.Content[j+1]
				found = true
				break
			}
		}

		if !found {
			keyNode = &yaml.Node{Kind: yaml.ScalarNode, Value: part}
			valueNode = &yaml.Node{Kind: yaml.MappingNode}

			if i == len(parts)-1 {
				valueNode = &yaml.Node{Kind: yaml.ScalarNode, Value: ToString(value)}
				if comment != "" {
					keyNode.LineComment = comment
				}
			}

			cur.Content = append(cur.Content, keyNode, valueNode)
		}

		cur = valueNode
	}
}

func yamlDescription(n *node) string {
	return n.Description
}

func yamlValue(n *node) string {
	if n.isSecret {
		return "<secret>"
	}

	return ToString(n.Value)
}
