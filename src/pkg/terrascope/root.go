package terrascope

type root struct {
	filename     string
	name         string
	ScopeTypes   []string          `hcl:"scopeTypes"`
	Dependencies []*rootDependency `hcl:"dependency,block"`
	ScopeMatches []*scopeMatch     `hcl:"scopeMatch,block"`
}

type rootDependency struct {
	RootName string            `hcl:"root"`
	Scopes   map[string]string `hcl:"scopes,optional"`
}

type scopeMatch struct {
	ScopeTypes map[string]string `hcl:"scopeTypes"`
}
