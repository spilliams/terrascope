package terrascope

import (
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

type buildContext struct {
	root  *root
	scope *CompiledScope
	*logrus.Entry
}

func newBuildContext(root *root, scope *CompiledScope, logger *logrus.Logger) *buildContext {
	return &buildContext{
		root:  root,
		scope: scope,
		Entry: logger.WithFields(logrus.Fields{
			"prefix": "builder",
			"root":   root.name,
			"scope":  scope.Address(),
		}),
	}
}

type rootConfig struct {
	Root       *root                 `hcl:"root,block"`
	Generators []*generator          `hcl:"generate,block"`
	Inputs     map[string]*cty.Value `hcl:"inputs,optional"`
}

type generator struct {
	ID       string `hcl:"id,label"`
	Path     string `hcl:"path,attr"`
	Contents string `hcl:"contents,attr"`
}

func (bc *buildContext) rootDirectory() string {
	return path.Dir(bc.root.filename)
}

func (bc *buildContext) destination() string {
	parts := []string{bc.rootDirectory(), ".terrascope"}
	parts = append(parts, bc.scope.ScopeValues...)
	return path.Join(parts...)
}

func BuildContext(bc *buildContext) (string, error) {
	rootVariable := cty.MapVal(map[string]cty.Value{
		"name": cty.StringVal(bc.root.name),
	})
	scopeVariable := bc.scope.ToCtyValue()
	attributesVariable := cty.ObjectVal(bc.scope.Attributes)

	bc.Trace("Building root")
	// first gotta reparse the config
	ctx := hclhelp.DefaultContext()
	ctx.Variables["root"] = rootVariable
	ctx.Variables["scope"] = scopeVariable
	ctx.Variables["attributes"] = cty.ObjectVal(bc.scope.Attributes)

	cfg := &rootConfig{}
	err := hclsimple.DecodeFile(bc.root.filename, ctx, cfg)
	if err != nil {
		return "", err
	}
	bc.Tracef("  fully decoded root config: %+v", cfg)

	destination := bc.destination()
	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return "", err
	}

	// TODO: empty the directory
	err = bc.copyAllFiles(bc.rootDirectory(), destination)
	if err != nil {
		return "", err
	}

	err = bc.processGenerators(cfg.Generators, destination)
	if err != nil {
		return "", err
	}

	err = bc.processInputs(cfg.Inputs, destination)
	if err != nil {
		return "", err
	}

	err = bc.generateDebugFile(destination, rootVariable, scopeVariable, attributesVariable)
	if err != nil {
		return "", err
	}

	return bc.destination(), nil
}
func (bc *buildContext) copyAllFiles(srcDir, destDir string) error {
	bc.Tracef("Walking %s", srcDir)
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
		bc.Tracef("  src: %s, dest: %s, file: %s, delta: %s, destPath: %s", srcDir, destDir, filepath, delta, destPath)

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

func (bc *buildContext) processGenerators(generators []*generator, destination string) error {
	bc.Trace("Processing generators")
	for _, gen := range generators {
		err := bc.processGenerator(gen, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bc *buildContext) processGenerator(gen *generator, destination string) error {
	destPath := path.Join(destination, gen.Path)
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(gen.Contents)
	bc.Tracef("contents: %s", gen.Contents)
	return err
}

func (bc *buildContext) processInputs(inputs map[string]*cty.Value, destination string) error {
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

func (bc *buildContext) generateDebugFile(destination string, rootVar, scopeVar, attrVar cty.Value) error {
	debugFile := hclwrite.NewEmptyFile()
	body := debugFile.Body()

	body.SetAttributeValue("root", rootVar)
	body.SetAttributeValue("scope", scopeVar)
	body.SetAttributeValue("attributes", attrVar)

	file, err := os.Create(path.Join(destination, ".terrascope.context.hcl"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = debugFile.WriteTo(file)
	return err
}
