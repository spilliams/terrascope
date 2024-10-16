package terrascope

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	hclhelp "github.com/spilliams/terrascope/internal/hcl"
	"github.com/zclconf/go-cty/cty"
)

type rootScopeContext struct {
	root  *Root
	scope *CompiledScope
	*logrus.Entry
}

func newRootScopeContext(root *Root, scope *CompiledScope, logger *logrus.Logger) *rootScopeContext {
	return &rootScopeContext{
		root:  root,
		scope: scope,
		Entry: logger.WithFields(logrus.Fields{
			"prefix": "builder",
			"root":   root.Name,
			"scope":  scope.Address(),
		}),
	}
}

func (rsc *rootScopeContext) String() string {
	return fmt.Sprintf("%s (%s)", rsc.root.Name, rsc.scope.Address())
}

type rootConfig struct {
	Root       *Root                 `hcl:"root,block"`
	Generators []*generator          `hcl:"generate,block"`
	Inputs     map[string]*cty.Value `hcl:"inputs,optional"`
}

type generator struct {
	ID       string `hcl:"id,label"`
	Path     string `hcl:"path,attr"`
	Contents string `hcl:"contents,attr"`
}

func (rsc *rootScopeContext) rootDirectory() string {
	return path.Dir(rsc.root.Filename)
}

func (rsc *rootScopeContext) destination() string {
	parts := []string{rsc.rootDirectory(), ".terrascope"}
	parts = append(parts, rsc.scope.ScopeValues...)
	return path.Join(parts...)
}

func BuildContext(rsc *rootScopeContext) (string, error) {
	rootVariable := cty.MapVal(map[string]cty.Value{
		"name": cty.StringVal(rsc.root.Name),
	})
	scopeVariable := rsc.scope.ToCtyValue()

	rsc.Trace("Building root")
	// first gotta reparse the config
	ctx := hclhelp.DefaultContext()
	ctx.Variables["root"] = rootVariable
	ctx.Variables["scope"] = scopeVariable
	ctx.Variables["attributes"] = cty.ObjectVal(rsc.scope.Attributes)

	cfg := &rootConfig{}
	err := hclsimple.DecodeFile(rsc.root.Filename, ctx, cfg)
	if err != nil {
		return "", err
	}
	rsc.Tracef("  fully decoded root config: %+v", cfg)

	destination := rsc.destination()

	err = os.RemoveAll(destination)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return "", err
	}

	err = rsc.copyAllFiles(rsc.rootDirectory(), destination)
	if err != nil {
		return "", err
	}

	err = rsc.processGenerators(cfg.Generators, destination)
	if err != nil {
		return "", err
	}

	err = rsc.processInputs(cfg.Inputs, destination)
	if err != nil {
		return "", err
	}

	err = rsc.generateDebugFile(destination, rootVariable, scopeVariable, rsc.scope.Attributes, rsc.scope.attributeSources)
	if err != nil {
		return "", err
	}

	return destination, nil
}

func (rsc *rootScopeContext) copyAllFiles(srcDir, destDir string) error {
	rsc.Tracef("Walking %s", srcDir)
	return filepath.WalkDir(srcDir, func(filepath string, d fs.DirEntry, err error) error {
		basename := path.Base(filepath)
		if d.IsDir() {
			if filepath == srcDir {
				// we don't want to skip the top directory, but we also don't
				// need to do anything with it
				return nil
			}
			if basename == ".terrascope" {
				return fs.SkipDir
			}
		}
		if basename == "terrascope.hcl" {
			return nil
		}

		// TODO: handle folders inside the source
		delta := strings.TrimPrefix(filepath, srcDir)
		destPath := path.Join(destDir, delta)
		rsc.Tracef("  src: %s, dest: %s, file: %s, delta: %s, destPath: %s", srcDir, destDir, filepath, delta, destPath)

		srcFile, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

func (rsc *rootScopeContext) processGenerators(generators []*generator, destination string) error {
	rsc.Trace("Processing generators")
	for _, gen := range generators {
		err := rsc.processGenerator(gen, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rsc *rootScopeContext) processGenerator(gen *generator, destination string) error {
	destPath := path.Join(destination, gen.Path)
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(gen.Contents)
	rsc.Tracef("contents: %s", gen.Contents)
	return err
}

func (rsc *rootScopeContext) processInputs(inputs map[string]*cty.Value, destination string) error {
	varsFile := hclwrite.NewEmptyFile()
	body := varsFile.Body()

	for k, v := range inputs {
		body.SetAttributeValue(k, *v)
	}

	file, err := os.Create(path.Join(destination, "terrascope.auto.tfvars"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = varsFile.WriteTo(file)
	return err
}

func (rsc *rootScopeContext) generateDebugFile(destination string, rootVar, scopeVar cty.Value, attrs map[string]cty.Value, attrSources map[string]string) error {
	debugFile := hclwrite.NewEmptyFile()
	body := debugFile.Body()

	body.SetAttributeValue("root", rootVar)
	body.SetAttributeValue("scope", scopeVar)
	body.SetAttributeRaw("attributes", hclhelp.TokensForObjectWithComments(attrs, attrSources))

	file, err := os.Create(path.Join(destination, ".terrascope.context.hcl"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = debugFile.WriteTo(file)
	return err
}

func CleanContext(rsc *rootScopeContext) (string, error) {
	destination := rsc.destination()
	err := os.RemoveAll(destination)
	return destination, err
}
