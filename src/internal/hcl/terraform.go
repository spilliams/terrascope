package hcl

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/sirupsen/logrus"
)

const separator = "."

// configFileSchema is the schema for the top-level of a config file. We use
// the low-level HCL API for this level so we can easily deal with each
// block type separately with its own decoding logic.
var configFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "terraform"},
		{Type: "provider", LabelNames: []string{"name"}},
		{Type: "variable", LabelNames: []string{"name"}},
		{Type: "locals"},
		{Type: "output", LabelNames: []string{"name"}},
		{Type: "module", LabelNames: []string{"name"}},
		{Type: "resource", LabelNames: []string{"id", "name"}},
		{Type: "data", LabelNames: []string{"type", "name"}},
	},
}

type configuration struct {
	locals map[string]*hcl.Attribute
	blocks map[string]*hcl.Block
	// ctx    *hcl.EvalContext
}

func newConfiguration() *configuration {
	return &configuration{
		locals: make(map[string]*hcl.Attribute, 0),
		blocks: make(map[string]*hcl.Block, 0),
		// ctx: &hcl.EvalContext{
		// 	Variables: make(map[string]cty.Value, 0),
		// 	Functions: map[string]function.Function{
		// 		"upper":  stdlib.UpperFunc,
		// 		"lower":  stdlib.LowerFunc,
		// 		"min":    stdlib.MinFunc,
		// 		"max":    stdlib.MaxFunc,
		// 		"strlen": stdlib.StrlenFunc,
		// 		"substr": stdlib.SubstrFunc,
		// 	},
		// },
	}
}

// Module is something that can parse a terraform module (configuration
// directory)
type Module interface {
	Configuration() *configuration
	Module() *tfconfig.Module
	ParseModuleDirectory(string) error
	Parser() *hclparse.Parser
	ParseTerraformFile(string) error
	DependencyGraph() (map[string][]string, error)
}

type module struct {
	cfg         *configuration
	fundamental *hclparse.Parser
	module      *tfconfig.Module
	*logrus.Logger
}

// NewModule builds an object that conforms to the ModuleParser interface
func NewModule(logger *logrus.Logger) Module {
	return &module{
		cfg:         newConfiguration(),
		fundamental: hclparse.NewParser(),
		Logger:      logger,
	}
}

func (m *module) ParseModuleDirectory(dirname string) error {
	module, diags := tfconfig.LoadModule(dirname)
	if diags.HasErrors() {
		return fmt.Errorf(diags.Error())
	}
	m.module = module

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		logrus.Debugf("examining file %s", f.Name())
		if strings.HasSuffix(f.Name(), ".tf") || strings.HasSuffix(f.Name(), ".hcl") {
			fullname := path.Join(dirname, f.Name())
			logrus.Debug("parsing!")
			if err := m.ParseTerraformFile(fullname); err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseTerraformFile parses a single file
func (m *module) ParseTerraformFile(filename string) error {
	file, diags := m.fundamental.ParseHCLFile(filename)
	if err := handleDiags(diags, m.fundamental.Files(), m.Logger.WriterLevel(logrus.WarnLevel)); err != nil {
		return err
	}

	content, _, diags := file.Body.PartialContent(configFileSchema)
	if err := handleDiags(diags, m.fundamental.Files(), m.Logger.WriterLevel(logrus.WarnLevel)); err != nil {
		return err
	}

	// read the contents into the receiver
	for _, block := range content.Blocks {
		var diags hcl.Diagnostics
		if block.Type == "locals" {
			var locals hcl.Attributes
			locals, diags = block.Body.JustAttributes()
			if err := handleDiags(diags, m.fundamental.Files(), m.Logger.WriterLevel(logrus.WarnLevel)); err != nil {
				return err
			}
			for name, attr := range locals {
				m.cfg.locals["local"+separator+name] = attr
			}
		} else {
			var blockName string
			if block.Type == "variable" {
				blockName = "var" + separator
			}
			if block.Type == "module" {
				blockName = "module" + separator
			}
			if block.Type == "output" {
				blockName = "output" + separator
			}
			if block.Type == "data" {
				blockName = "data" + separator
			}
			blockName += strings.Join(block.Labels, separator)
			m.cfg.blocks[blockName] = block
		}
		if err := handleDiags(diags, m.fundamental.Files(), m.Logger.WriterLevel(logrus.WarnLevel)); err != nil {
			return err
		}
	}

	return nil
}

func (m *module) DependencyGraph() (map[string][]string, error) {
	graph := make(map[string][]string, 0)

	for name, local := range m.cfg.locals {
		graph[name] = attributeDependencies(local)
	}

	for name, block := range m.cfg.blocks {
		deps, err := blockDependencies(block)
		if err != nil {
			return nil, err
		}
		graph[name] = deps
	}

	m.Debugf("dependency graph before pruning for index and unknowns: %#v", graph)

	// dependencies go all the way through attribute names. For instance, an
	// output.stack_name might depend on random_string.slug.result, but the
	// m.cfg.blocks map only includes a "random_string.slug".
	// So for each dependency we need to check: is it in our cfg? Or do we need
	// to truncate it?
	// Also, the downstreams here include things like "string" (the raw type
	// used by a variable's "type" attribute) or "path.module", which is a
	// Terraform builtin. So we only want to include things we recognize
	for upstream, deps := range graph {
		keepers := make([]string, 0, len(deps))
		for _, downstream := range deps {
			parts := strings.Split(downstream, separator)
			for limit := len(parts); limit > 0; limit-- {
				trial := strings.Join(parts[0:limit], separator)
				if !m.Has(trial) {
					continue
				}
				keepers = append(keepers, trial)
				break
			}
		}

		// and the resulting map might have dupes in the values, so uniq them.
		graph[upstream] = unique(keepers)
	}

	return graph, nil
}

func unique[T comparable](list []T) []T {
	uniq := make([]T, 0, len(list))
	truth := make(map[T]bool)

	for _, val := range list {
		if _, ok := truth[val]; !ok {
			truth[val] = true
			uniq = append(uniq, val)
		}
	}
	return uniq
}

func (m *module) Has(path string) bool {
	if m.cfg.locals[path] != nil {
		return true
	}
	if m.cfg.blocks[path] != nil {
		return true
	}
	return false
}

func attributeDependencies(attr *hcl.Attribute) []string {
	return expressionDependencies(attr.Expr)
}

func blockDependencies(block *hcl.Block) ([]string, error) {
	deps := make([]string, 0)

	attrs, _ := block.Body.JustAttributes()
	for _, attr := range attrs {
		deps = append(deps, attributeDependencies(attr)...)
	}

	// TODO how can we take into account exprs of nested blocks?
	return deps, nil
}

func expressionDependencies(expr hcl.Expression) []string {
	deps := make([]string, 0)
	for _, traversal := range expr.Variables() {
		var varName string
		for i, step := range traversal {
			if i == 0 {
				varName = step.(hcl.TraverseRoot).Name
				continue
			}
			varName += separator + step.(hcl.TraverseAttr).Name
		}
		deps = append(deps, varName)
	}

	return deps
}

func (m *module) Parser() *hclparse.Parser {
	return m.fundamental
}

func (m *module) Module() *tfconfig.Module {
	return m.module
}

func (m *module) Configuration() *configuration {
	return m.cfg
}
