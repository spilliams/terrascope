package terrascope

import (
	"fmt"
	"path"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spilliams/terrascope/internal/generate"
	hclhelp "github.com/spilliams/terrascope/internal/hcl"
)

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

// ParseRoot tells the receiver to parse a root module configuration file at the
// given path.
func (p *Project) ParseRoot(cfgFile string) (*root, error) {
	// partial decode first, because we don't know what scope or attributes
	// this config will use. We're just looking for the `root` block here.
	cfg := &struct {
		Root *root `hcl:"root,block"`
	}{}

	err := hclsimple.DecodeFile(cfgFile, hclhelp.DefaultContext(), cfg)
	r := cfg.Root
	// we purposefully ignore err until the end
	if r == nil {
		p.Warnf("Root detected at %s failed to decode. Does it have a complete terrascope.hcl file?", cfgFile)
		return nil, nil
	}

	r.filename = cfgFile
	r.name = path.Base(path.Dir(cfgFile))
	return r, err
}

// GenerateRoot creates and runs a new root generator, using the receiver's
// scope types.
func (p *Project) GenerateRoot(name string) error {
	if len(p.ScopeTypes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terrascope `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.ScopeTypes))
	for i, el := range p.ScopeTypes {
		scopeTypes[i] = el.Name
	}

	return generate.Root(name, p.RootsDir, scopeTypes, p.Logger)
}

func (p *Project) assertRootDependenciesAcyclic() error {
	visited := make(map[string]bool)
	for rootName := range p.Roots {
		if visited[rootName] {
			continue
		}
		stack := make(map[string]bool)
		if has, list := hasCyclicDependency(rootName, p.Roots, visited, stack); has {
			return fmt.Errorf("cyclical dependency detected for root '%s': %+v", rootName, list)
		}
	}
	return nil
}

// hasCyclicDependency performs a depth-first search. It takes the current root
// name, the map of all roots, a visited map to keep track of which nodes have
// been visited, and a stack map to keep track of nodes we have yet to visit.
// It returns true if a cyclical dependency is found, and false otherwise.
// When it returns true, it also returns the list of root names in the cycle.
func hasCyclicDependency(rootName string, roots map[string]*root, visited, stack map[string]bool) (bool, []string) {
	visited[rootName] = true
	stack[rootName] = true

	for _, dep := range roots[rootName].Dependencies {
		if !visited[dep.RootName] {
			if has, _ := hasCyclicDependency(dep.RootName, roots, visited, stack); has {
				return true, keys(visited)
			}
		} else if stack[dep.RootName] {
			return true, keys(visited)
		}
	}

	delete(stack, rootName)
	return false, []string{}
}

func keys(m map[string]bool) []string {
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}
	return l
}
