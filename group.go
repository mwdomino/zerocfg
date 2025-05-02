package zerocfg

// Grp represents a group of configuration options that can share a common prefix and set of option modifiers.
// Groups are useful for organizing related configuration options and applying shared modifiers.
type Grp struct {
	prefix string
	opts   []OptNode
}

// NewGroup creates a new Grp with the specified prefix and option modifiers.
// All options added to this group will have the prefix prepended to their names.
//
// Example:
//
//	g := NewGroup("db")
//	Str("host", "localhost", "db host", Group(g)) // becomes "db.host"
func NewGroup(prefix string, opts ...OptNode) *Grp {
	return &Grp{
		prefix: prefix,
		opts:   opts,
	}
}

// NewOptions creates a new Grp with no prefix but with shared option modifiers.
// This is useful for applying the same modifiers (e.g., Secret, Required) to multiple options without a hierarchical prefix.
//
// Example:
//
//	g := NewOptions(Secret())
//	Str("api_key", "", "API key", Group(g)) // marked as secret
func NewOptions(opts ...OptNode) *Grp {
	return &Grp{
		opts: opts,
	}
}

func (g *Grp) key(child string) string {
	if g.prefix == "" {
		return child
	}

	return g.prefix + "." + child
}

func (g *Grp) applyOpts(n *Node) {
	for _, opt := range g.opts {
		opt(n)
	}
}
