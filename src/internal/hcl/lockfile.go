package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// lockfileSchema is the hcl schema we expect a lockfile to conform to
var lockfileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "provider", LabelNames: []string{"id"}},
	},
}

// LockfileProvider represents a single `provider` block of a Lockfile
type LockfileProvider struct {
	ID          string   `hcl:"id,label"`
	Version     string   `hcl:"version"`
	Constraints string   `hcl:"constraints,optional"`
	Hashes      []string `hcl:"hashes,optional"`
}

// Lockfile represents a Terraform lockfile
type Lockfile struct {
	Providers []*LockfileProvider `hcl:"provider,block"`
}

// CompactProviders returns a list of the receiver's providers, formatted in
// short (compact) strings, e.g. "name@version".
func (lf *Lockfile) CompactProviders() []string {
	providers := make([]string, len(lf.Providers))
	for i, provider := range lf.Providers {
		providers[i] = fmt.Sprintf("%s@%s", provider.ID, provider.Version)
	}
	return providers
}

func (lp *LockfileProvider) String() string {
	return fmt.Sprintf("<Provider: %s, version %s; constraint %s; %d hashes>", lp.ID, lp.Version, lp.Constraints, len(lp.Hashes))
}

// ParseLockfile parses a given file as a Terraform lockfile, and returns a
// Lockfile object representing the configuration.
// Will return an error if: the HCL parser fails to parse the file, or fails to
// read the file into the lockfile schema, or the file contains hcl blocks other
// than providers.
func ParseLockfile(filename string) (*Lockfile, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}

	lockfile := &Lockfile{}
	lockfile.Providers = make([]*LockfileProvider, 0)
	content, diags := f.Body.Content(lockfileSchema)
	if err := handleDiags(diags, parser.Files(), nil); err != nil {
		return nil, err
	}
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: map[string]function.Function{},
	}
	for _, block := range content.Blocks {
		if block.Type != "provider" {
			return nil, fmt.Errorf("found unknown block type %s: %+v", block.Type, block)
		}
		provider := LockfileProvider{}
		diags := gohcl.DecodeBody(block.Body, ctx, &provider)
		if err := handleDiags(diags, parser.Files(), nil); err != nil {
			return nil, err
		}
		provider.ID = block.Labels[0]
		lockfile.Providers = append(lockfile.Providers, &provider)
	}

	return lockfile, nil
}
