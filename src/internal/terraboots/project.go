package terraboots

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/internal/scopedata"
)

// ProjectConfig represents the configuration file of a Terraboots project
type ProjectConfig struct {
	Project *Project `hcl:"terraboots,block"`
}

// Project represents a single Terraboots project, complete with scope types,
// scope data, and root configurations.
type Project struct {
	configFile string
	ID         string `hcl:"id,label"`

	ScopeTypes     []*ScopeType `hcl:"scope,block"`
	ScopeDataFiles []string     `hcl:"scopeData"`
	compiledScopes scopedata.CompiledScopes

	RootsDir string `hcl:"rootsDir"`
	Roots    map[string]*root

	*logrus.Entry
}

// ScopeType represents a single scope available to a project
type ScopeType struct {
	Name         string `hcl:"name"`
	Description  string `hcl:"description,optional"`
	DefaultValue string `hcl:"default,optional"`
	// Validations  []*ProjectScopeValidation `hcl:"validation,block"`
}

func (sc *ScopeType) String() string {
	return sc.Name
}

// ProjectScopeValidation
// type ProjectScopeValidation struct {
// 	Condition    bool   `hcl:"condition"`
// 	ErrorMessage string `hcl:"error_message"`
// }

// ParseProject reads the given configuration file and parses it as a new
// Terraboots project
func ParseProject(cfgFile string, logger *logrus.Logger) (*Project, error) {
	cfg := &ProjectConfig{}
	cfgFile, err := filepath.Abs(cfgFile)
	if err != nil {
		return nil, err
	}
	err = hclsimple.DecodeFile(cfgFile, nil, cfg)
	if err != nil {
		return nil, err
	}

	project := cfg.Project
	project.configFile = cfgFile
	project.Entry = logger.WithField("prefix", "project")

	err = project.readScopeData()
	if err != nil {
		return nil, err
	}

	project.Debugf("Project has %d compiled scopes", project.compiledScopes.Len())
	for _, scope := range project.compiledScopes {
		project.Trace(scope.Address())
	}
	return project, nil
}

func (p *Project) projectDir() string {
	return path.Dir(p.configFile)
}

// BuildRoot tells the receiver to build a root module, and returns a list of
// directories where the root was built to.
func (p *Project) BuildRoot(rootName string, scopes []string) ([]string, error) {
	// first, get the root
	root, ok := p.Roots[rootName]
	if !ok {
		var err error
		root, err = p.AddRoot(rootName)
		if err != nil {
			return nil, err
		}
	}
	p.Debugf("root: %+v", root)

	// what scopes to build for?
	matchingScopes, err := p.determineMatchingScopes(root, scopes)
	if err != nil {
		return nil, err
	}
	p.Infof("Root will be built for %d %s", len(matchingScopes), pluralize("scope", "scopes", len(matchingScopes)))
	for _, scope := range matchingScopes {
		p.Trace(scope.Address())
	}

	builds := make([]*buildContext, len(matchingScopes))
	for i, scope := range matchingScopes {
		builds[i] = newBuildContext(root, scope, p.Entry.Logger)
	}

	// TODO: root dependencies. Do the buildContexts figure it out? Then we need
	// to phase them and deduplicate them

	// TODO: use a worker pool
	dirs := make([]string, len(builds))
	for i, build := range builds {
		err := build.Build()
		if err != nil {
			return nil, err
		}

		dirs[i] = build.destination()
	}

	return dirs, nil
}

// AddRoot tells the receiver to add a root module to its internal records.
// The `rootName` must be a directory name located in the receiver's `RootsDir`.
func (p *Project) AddRoot(rootName string) (*root, error) {
	// look for named root
	rootDir := path.Join(p.RootsDir, rootName)
	_, err := os.Stat(rootDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("could not locate a root named '%s' in the roots directory '%s'", rootName, p.RootsDir)
	} else if err != nil {
		return nil, err
	}

	// look for terraboots file
	rootCfg := path.Join(rootDir, "terraboots.hcl")
	_, err = os.Stat(rootCfg)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("found a root named '%s' in the roots directory '%s', but it does not contain a terraboots.hcl configuration", rootName, p.RootsDir)
	} else if err != nil {
		return nil, err
	}

	r, err := ParseRoot(rootCfg)
	if r == nil && err != nil {
		return nil, err
	}

	if p.Roots == nil {
		p.Roots = make(map[string]*root)
	}
	p.Roots[r.ID] = r
	return r, nil
}

// determineMatchingScopes takes in a root configuration and an optional list of
// scopes. It returns a list of scopedata.CompiledScopes where each scope in the
// list (a) matches at least one scopeMatch expression of the root, and
// (b) matches at least one scope given.
// Note that a root with no scopeMatch expressions will be treated as if all its
// scope types allow all values (`.*`).
func (p *Project) determineMatchingScopes(root *root, scopes []string) (scopedata.CompiledScopes, error) {
	matchingScopes := scopedata.CompiledScopes{}
	// if they don't specify any scope matches, assume .* for all
	if root.ScopeMatches == nil || len(root.ScopeMatches) == 0 {
		allScopeMatchTypes := make(map[string]string)
		for _, scope := range root.ScopeTypes {
			allScopeMatchTypes[scope] = ".*"
		}
		root.ScopeMatches = []*scopeMatch{
			{ScopeTypes: allScopeMatchTypes},
		}
	}

	for _, scopeMatch := range root.ScopeMatches {
		matches, err := p.compiledScopes.Matching(scopeMatch.ScopeTypes)
		if err != nil {
			return nil, err
		}
		matchingScopes = append(matchingScopes, matches...)
	}
	matchingScopes = matchingScopes.Deduplicate()
	sort.Sort(matchingScopes)

	// also abide by this list
	if len(scopes) > 0 {
		filteredMatchingScopes := scopedata.CompiledScopes{}
		scopeFilters := make([]map[string]string, len(scopes))
		for i, scope := range scopes {
			scopeFilter, err := p.makeScopeFilter(scope)
			if err != nil {
				return nil, err
			}
			scopeFilters[i] = scopeFilter
		}
		p.Debugf("filters on the root's full list of scope values:\n%+v", scopeFilters)
		for _, scope := range matchingScopes {
			for _, filter := range scopeFilters {
				ok, err := scope.Matches(filter)
				if err != nil {
					return nil, err
				}
				if ok {
					filteredMatchingScopes = append(filteredMatchingScopes, scope)
				}
			}
		}
		matchingScopes = filteredMatchingScopes
	}

	if len(matchingScopes) == 0 {
		return nil, fmt.Errorf("No matching scope values found.\n"+
			"Root '%s' applies to the scope types %v.\n"+
			"All scopes in the project were searched (%d), and none matched these types. Please provide\n"+
			"\t- new scope data for the project,\n"+
			"\t- different scope types in the root configuration file, or\n"+
			"\t- new scope matches in the root configuration file.",
			root.ID, root.ScopeTypes, len(p.compiledScopes))
	}
	return matchingScopes, nil
}

func pluralize(single, plural string, count int) string {
	if count == 1 {
		return single
	}
	return plural
}
