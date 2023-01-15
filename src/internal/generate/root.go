package generate

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

func Root(name string, dir string, scopeTypes []string, logger *logrus.Logger) error {
	rg := &rootGenerator{
		name:       name,
		dir:        dir,
		scopeTypes: scopeTypes,
		Entry:      logger.WithField("prefix", "rootgen"),
	}
	return rg.Run()
}

type rootGenerator struct {
	name       string
	dir        string
	scopeTypes []string
	*logrus.Entry
}

func (rg *rootGenerator) Run() error {
	root, err := rg.surveyForRootConfiguration()
	if err != nil {
		return err
	}
	rg.Debugf("root: %#v", root)

	hclfile := rg.generateRootConfigurationFile(root)
	rg.Debugf("root hcl:\n%s", hclfile.Bytes())
	return rg.writeRootConfigurationFile(hclfile.Bytes())
}

type root struct {
	Name         string
	ScopeType    string
	scopeMatches []*scopeMatch
}

type scopeMatch struct {
	ScopeTypes map[string]string
}

func (sm *scopeMatch) WriteAnswer(field string, value interface{}) error {
	if sm.ScopeTypes == nil {
		sm.ScopeTypes = make(map[string]string)
	}

	sValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	sm.ScopeTypes[field] = sValue
	return nil
}

func (rg *rootGenerator) surveyForRootConfiguration() (*root, error) {
	questions := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What is the root's name?",
				Default: rg.name,
			},
			// TODO Validate: hcl label
		},
		{
			Name: "scopeType",
			Prompt: &survey.Select{
				Message: "What scope type should the root be built for?",
				Options: rg.scopeTypes,
			},
		},
	}
	// TODO: ask for dependencies. Should we be allowed to ask for any other
	// module in the project? or only those who are compatible with our answer
	// for scopeType? (We probably don't know the answer for scopeType here, nor
	// do we know it during the question Validator, it'll have to be afterwards)

	cfg := root{}
	err := survey.Ask(questions, &cfg)
	cfg.scopeMatches, err = rg.surveyForRootScopeMatches(cfg.ScopeType)
	return &cfg, err
}

func (rg *rootGenerator) surveyForRootScopeMatches(rootScopeType string) ([]*scopeMatch, error) {
	scopeMatches := make([]*scopeMatch, 0)
	isSubsequent := false
	rootScopeTypes := make([]string, 0, len(rg.scopeTypes))
	allValue := make(map[string]string)
	for _, scopeType := range rg.scopeTypes {
		rootScopeTypes = append(rootScopeTypes, scopeType)
		allValue[scopeType] = ".*"
		if scopeType == rootScopeType {
			break
		}
	}
	allOption := fmt.Sprintf("%s.* (all)", strings.Join(rootScopeTypes, ".*."))

	customOption := "custom"
	doneOption := "done"

	for true {
		var value string

		message := "What scope values should the root be built for?"
		options := []string{
			allOption,
			customOption,
		}
		if isSubsequent {
			message = "Any additional scope values the root should be built for?"
			options = append(options, doneOption)
		}
		isSubsequent = true

		err := survey.AskOne(&survey.Select{
			Message: message,
			Options: options,
		}, &value)
		if err == terminal.InterruptErr {
			return nil, err
		}
		if err != nil {
			rg.Warn(err)
		}

		if value == doneOption {
			break
		}

		if value == allOption {
			scopeMatches = append(scopeMatches, &scopeMatch{allValue})
			continue
		}

		if value == customOption {
			customMatch, err := surveyForCustomScopeMatch(rootScopeTypes)
			if err != nil {
				rg.Warn(err)
				continue
			}
			scopeMatches = append(scopeMatches, customMatch)
		}
	}
	return scopeMatches, nil
}

func surveyForCustomScopeMatch(scopeTypes []string) (*scopeMatch, error) {
	questions := make([]*survey.Question, len(scopeTypes))
	for i, scopeType := range scopeTypes {
		question := &survey.Question{
			Name:   scopeType,
			Prompt: &survey.Input{Message: fmt.Sprintf("What values are allowed for the scope type %s? You may use regular expressions", scopeType)},
		}
		questions[i] = question
	}

	fmt.Printf("%#v\n", questions)
	answers := scopeMatch{}
	err := survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}

	return &answers, nil
}

func (rg *rootGenerator) generateRootConfigurationFile(cfg *root) *hclwrite.File {
	f := hclwrite.NewEmptyFile()
	cfgBody := f.Body()

	rootBlock := cfgBody.AppendNewBlock("root", []string{cfg.Name})
	rootBody := rootBlock.Body()

	scopeTypeVals := make([]cty.Value, 0, len(rg.scopeTypes))
	for _, scopeType := range rg.scopeTypes {
		scopeTypeVals = append(scopeTypeVals, cty.StringVal(scopeType))
		if scopeType == cfg.ScopeType {
			break
		}
	}
	rootBody.SetAttributeRaw("scopeTypes", hclwrite.TokensForValue(cty.ListVal(scopeTypeVals)))

	for _, scopeMatch := range cfg.scopeMatches {
		matchBlock := rootBody.AppendNewBlock("scopeMatch", []string{})
		matchBody := matchBlock.Body()
		matchTypesVal := make(map[string]cty.Value)
		for matchType, scopeValue := range scopeMatch.ScopeTypes {
			matchTypesVal[matchType] = cty.StringVal(scopeValue)
		}
		matchBody.SetAttributeRaw("scopeTypes", hclwrite.TokensForValue(cty.MapVal(matchTypesVal)))
	}

	return f
}

func (rg *rootGenerator) writeRootConfigurationFile(b []byte) error {
	rootDir := path.Join(rg.dir, rg.name)
	err := os.MkdirAll(rootDir, 0755)
	if err != nil {
		return err
	}
	rootFile := path.Join(rootDir, "terraboots.hcl")
	file, err := os.Create(rootFile)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	rg.Infof("New root configuration file %s created.", rootFile)
	return err
}
