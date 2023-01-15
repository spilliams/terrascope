package projectgen

import (
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/internal/surveyhelp"
	"github.com/zclconf/go-cty/cty"
)

type projectConfiguration struct {
	ProjectName string
	RootDir     string
	ScopeData   string
	ScopeTypes  []string
}

func GenerateProject(logger *logrus.Logger) error {
	g := &generator{
		Entry: logger.WithField("prefix", "projectgen"),
	}
	return g.run()
}

type generator struct {
	*logrus.Entry
}

func (g *generator) run() error {
	questions := []*survey.Question{
		{
			Name:     "projectName",
			Prompt:   &survey.Input{Message: "What is the name of the project?"},
			Validate: survey.Required,
			// TODO: validate hcl label (no spaces?)
		},
		{
			Name: "rootDir",
			Prompt: &survey.Input{
				Message: "Where do you want to store your roots?",
				Default: "terraform/roots",
			},
			Validate: survey.Required,
			// TODO: validate ok for a path
		},
		{
			Name: "scopeData",
			Prompt: &survey.Input{
				Message: "Where do you want to store your scope data?",
				Default: "data.hcl",
			},
			Validate: survey.Required,
			// TODO: validate ok for a file
		},
		{
			Name:      "scopeTypes",
			Prompt:    &survey.Input{Message: "What scope types does the project use? (in order, space-delimited)"},
			Transform: surveyhelp.SplitTransformer,
		},
	}

	answers := projectConfiguration{}
	err := survey.Ask(questions, &answers)
	if err != nil {
		return err
	}

	g.Debugf("Answers: %+v", answers)
	if err := g.writeProjectTerrabootsFile(answers); err != nil {
		return err
	}

	return g.createRootsDirectory(answers.RootDir)

}

func (g *generator) writeProjectTerrabootsFile(cfg projectConfiguration) error {
	f := hclwrite.NewEmptyFile()
	projectBody := f.Body()

	tbBlock := projectBody.AppendNewBlock("terraboots", []string{cfg.ProjectName})
	tbBody := tbBlock.Body()
	tbBody.SetAttributeRaw("rootsDir", hclwrite.TokensForValue(cty.StringVal(cfg.RootDir)))
	tbBody.SetAttributeRaw("scopeData", hclwrite.TokensForValue(cty.ListVal([]cty.Value{cty.StringVal(cfg.ScopeData)})))

	for _, scope := range cfg.ScopeTypes {
		scopeBlock := tbBody.AppendNewBlock("scope", []string{})
		scopeBody := scopeBlock.Body()
		scopeBody.SetAttributeRaw("name", hclwrite.TokensForValue(cty.StringVal(scope)))
	}

	filename := "terraboots.hcl"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = file.Write(f.Bytes())
	if err != nil {
		return err
	}
	g.Infof("New project configuration file %s created.", filename)
	return nil
}

func (g *generator) createRootsDirectory(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(path.Join(dir, ".gitkeep"))
	defer file.Close()
	if err != nil {
		return err
	}
	g.Infof("New root module directory %s created.", dir)
	return nil
}
