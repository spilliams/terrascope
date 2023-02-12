package generate

import (
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/spilliams/terrascope/internal/surveyhelp"
	"github.com/zclconf/go-cty/cty"
)

type projectConfiguration struct {
	ProjectName string
	RootDir     string
	ScopeData   string
	ScopeTypes  []string
}

type projectGenerator struct {
	*logrus.Entry
}

func Project(logger *logrus.Logger) error {
	g := &projectGenerator{
		Entry: logger.WithField("prefix", "projectgen"),
	}
	return g.Run()
}

func (pg *projectGenerator) Run() error {
	answers, err := surveyForProjectConfiguration()
	if err != nil {
		return err
	}
	pg.Debugf("Answers: %+v", answers)

	hclfile := generateProjectConfigurationFile(answers)

	if err := pg.writeProjectTerrascopeFile(hclfile.Bytes()); err != nil {
		return err
	}

	return pg.createRootsDirectory(answers.RootDir)
}

func surveyForProjectConfiguration() (*projectConfiguration, error) {
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
	return &answers, err
}

func generateProjectConfigurationFile(cfg *projectConfiguration) *hclwrite.File {
	f := hclwrite.NewEmptyFile()
	projectBody := f.Body()

	tbBlock := projectBody.AppendNewBlock("terrascope", []string{cfg.ProjectName})
	tbBody := tbBlock.Body()
	tbBody.SetAttributeRaw("rootsDir", hclwrite.TokensForValue(cty.StringVal(cfg.RootDir)))
	tbBody.SetAttributeRaw("scopeData", hclwrite.TokensForValue(cty.ListVal([]cty.Value{cty.StringVal(cfg.ScopeData)})))

	for _, scope := range cfg.ScopeTypes {
		scopeBlock := tbBody.AppendNewBlock("scope", []string{})
		scopeBody := scopeBlock.Body()
		scopeBody.SetAttributeRaw("name", hclwrite.TokensForValue(cty.StringVal(scope)))
	}

	return f
}

func (pg *projectGenerator) writeProjectTerrascopeFile(b []byte) error {
	filename := "terrascope.hcl"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	pg.Infof("New project configuration file %s created.", filename)
	return nil
}

func (pg *projectGenerator) createRootsDirectory(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(path.Join(dir, ".gitkeep"))
	defer file.Close()
	if err != nil {
		return err
	}
	pg.Infof("New root module directory %s created.", dir)
	return nil
}
