package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const data = `scope "foo" "top" {
	account = "12345"
	
	scope "bar" "next" {
		vpc_cidr = "10.0.0.0/16"
	}
}`

func main() {
	cfg, err := parseSimple(data)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("%+v", cfg.Scopes)
}

type config struct {
	Scopes []*scope `hcl:"scope,block"`
}

type scope struct {
	Type     string         `hcl:"type,label"`
	Name     string         `hcl:"name,label"`
	Children []*scope       `hcl:"scope,block"`
	Attrs    hcl.Attributes `hcl:",remain"`
}

func (s *scope) String() string {
	attrs := []string{}
	for k := range s.Attrs {
		attrs = append(attrs, k)
	}
	return fmt.Sprintf("<scope.%s.%s: %d children, attrs: %s>", s.Type, s.Name, len(s.Children), attrs)
}

func parseSimple(data string) (*config, error) {
	cfg := &config{}
	err := hclsimple.Decode("main.hcl", []byte(data), nil, cfg)
	if err = handleSimpleDecodeError(err); err != nil {
		log.Println("error decoding scope data")
		return cfg, err
	}
	return cfg, nil
}

func handleSimpleDecodeError(err error) error {
	diags, typeOK := err.(hcl.Diagnostics)
	if !typeOK {
		return err
	}

	var newDiags hcl.Diagnostics
	for _, diag := range diags {
		if diag.Summary != "Unexpected \"scope\" block" {
			newDiags = append(newDiags, diag)
		}
	}

	if len(newDiags) > 0 {
		return newDiags
	}
	return nil
}

func parseSimpleDebug(data string) (*config, error) {
	cfg := &config{}
	filename := "main.hcl"
	src := []byte(data)

	var file *hcl.File
	var diags hcl.Diagnostics

	file, diags = hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	diags = gohcl.DecodeBody(file.Body, nil, cfg)
	if diags.HasErrors() {
		return nil, diags
	}
	return cfg, nil
}

func parseSchema(data string) (*config, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL([]byte(data), "main.go")
	if diags.HasErrors() {
		return nil, diags
	}

	scopeSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "scope", LabelNames: []string{"type", "id"}},
		},
	}
	content, _, diags := f.Body.PartialContent(scopeSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	cfg := &config{
		Scopes: make([]*scope, 0),
	}
	for _, block := range content.Blocks {
		content2, leftover, diags := block.Body.PartialContent(scopeSchema)
		if diags.HasErrors() {
			log.Println("could not get partial content from block")
			return nil, diags
		}
		log.Printf("content2: %+v", content2)
		log.Printf("  blocks: %+v", content2.Blocks)
		log.Printf("leftover: %+v", leftover)
		// ignore diags form this until I can parse them to just ignore the
		// "Unexpected %q block" type
		attrs, diags2 := leftover.JustAttributes()
		if diags2.HasErrors() {
			log.Println("could not get attrs from leftovers")
			return nil, diags2
		}
		for k, v := range attrs {
			value, diags3 := v.Expr.Value(nil)
			log.Printf("attr %s = %+v", k, value.AsString())
			if diags3.HasErrors() {
				log.Println("could not evaluate attribute expression")
				return nil, diags3
			}
		}

		scope := &scope{}
		diags = gohcl.DecodeBody(block.Body, nil, scope)
		if diags.HasErrors() {
			return nil, diags
		}
		cfg.Scopes = append(cfg.Scopes, scope)
	}

	return cfg, nil
}
