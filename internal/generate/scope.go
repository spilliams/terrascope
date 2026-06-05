package generate

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	hclhelp "github.com/spilliams/terrascope/internal/hcl"
)

// Scope generates a new scope data file with the given scope types. This
// function will survey the user for some necessary information, via stdin.
func Scope(scopeTypes []string, filename string, logger *logrus.Logger) error {
	sg := &scopeGenerator{
		filename:   filename,
		scopeTypes: scopeTypes,
		Entry:      logger.WithField("prefix", "scopegen"),
	}
	return sg.Run()
}

// scopeGenerator stores an ordered list of scope types, a filename to store data to,
// and composes a Logger for debugging
type scopeGenerator struct {
	filename   string
	scopeTypes []string
	*logrus.Entry
}

// Run surveys the user about their scope value, and returns bytes representing
// an hcl file
func (sg *scopeGenerator) Run() error {
	rootScopes, err := sg.surveyForScopeValues()
	if err != nil {
		return err
	}
	if len(rootScopes) == 0 {
		sg.Warn("No scopes were generated, exiting.")
		return nil
	}

	hclfile := generateScopeDataFile(rootScopes)

	return sg.writeScopeDataFile(hclfile.Bytes())
}

const answerRE = "[0-9a-zA-Z-_]"

var helpText = fmt.Sprintf("Answers must be space-separated, and may consist of the characters %s\n"+
	"Leave any answer blank to mark the current scope as complete with no children.\n"+
	"Press Ctrl+C at any time to cancel.", answerRE)

type nestedScope struct {
	Type     string `hcl:"type,label"`
	Name     string `hcl:"name,label"`
	Address  string
	Children []*nestedScope `hcl:"scope,block"`
	Attrs    hcl.Attributes `hcl:",remain"`

	scopeTypeIndex int
}

// surveyForScopeValues uses the receiver's scopeTypes to ask the user for all
// the different values of the scopes.
// Returns a list of the top-level scope values (the scope values for the first
// scope type)
func (sg *scopeGenerator) surveyForScopeValues() ([]*nestedScope, error) {
	sg.Infof("Scope types in this projct, in order, are: %s", strings.Join(sg.scopeTypes, ", "))

	// First one's free
	firstValues, err := askScope("What are the allowable scope values for `%s`?\n", sg.scopeTypes[0])
	if err != nil {
		return nil, err
	}
	if len(firstValues) == 0 {
		sg.Debugf("empty value, exiting")
		return nil, nil
	}
	if err := validateScope(firstValues); err != nil {
		return nil, err
	}
	sg.Debugf("read new scope values %v", firstValues)

	roots := make([]*nestedScope, len(firstValues))
	prompts := make([]*nestedScope, len(firstValues))

	for i, el := range firstValues {
		value := &nestedScope{
			Name:           el,
			Type:           sg.scopeTypes[0],
			scopeTypeIndex: 0,
			Children:       make([]*nestedScope, 0),
		}
		value.Address = fmt.Sprintf("%s.%s", value.Type, value.Name)
		roots[i] = value
		prompts[i] = value
	}

	for len(prompts) > 0 {
		prompt := prompts[0]
		prompts = prompts[1:]

		if prompt.scopeTypeIndex+1 == len(sg.scopeTypes) {
			// this scope value cannot have children
			continue
		}

		values, err := askScope("Within %s, what are the allowable scope values for `%s`?\n", prompt.Address, sg.scopeTypes[prompt.scopeTypeIndex+1])
		if err != nil {
			return nil, err
		}

		if len(values) == 0 {
			sg.Debugf("user entered no scope values for prompt, closing this scope")
			continue
		}

		if err := validateScope(values); err != nil {
			return nil, err
		}

		sg.Debugf("read new scope values %v", values)

		for _, el := range values {
			value := &nestedScope{
				Name:           el,
				Type:           sg.scopeTypes[prompt.scopeTypeIndex+1],
				scopeTypeIndex: prompt.scopeTypeIndex + 1,
				Children:       make([]*nestedScope, 0),
			}
			value.Address = strings.Join([]string{prompt.Address, string(value.Type), value.Name}, ".")
			prompt.Children = append(prompt.Children, value)
			prompts = append(prompts, value)
		}
	}

	sg.Debugf("%+v", roots)

	return roots, nil
}

func askScope(format string, a ...any) ([]string, error) {
	var value string
	err := survey.AskOne(&survey.Input{
		Message: fmt.Sprintf(format, a...),
		Help:    helpText,
	}, &value)
	if len(value) == 0 {
		return []string{}, err
	}
	return strings.Split(value, " "), err
}

func validateScope(answers []string) error {
	seen := make(map[string]bool)
	for _, answer := range answers {
		if len(answer) == 0 {
			return fmt.Errorf("Values cannot be blank (did you press space twice?)")
		}
		if _, ok := seen[answer]; ok {
			return fmt.Errorf("Cannot use the same scope value (%s) more than once for a single scope type", answer)
		}
		seen[answer] = true
		re := regexp.MustCompile(fmt.Sprintf("^%s+$", answerRE))
		if !re.MatchString(answer) {
			return fmt.Errorf("Scope value '%s' does not match valid character set %s", answer, answerRE)
		}
	}
	return nil
}

// generateScopeDataFile reads the given scopes and produces an `hclwrite.File`
// object that is ready to be written to disk.
func generateScopeDataFile(rootScopes []*nestedScope) *hclwrite.File {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	rootBody.AppendUnstructuredTokens(hclhelp.TokensForComment("This file was generated by terrascope"))
	rootBody.AppendNewline()

	for _, root := range rootScopes {
		rootBody = addScopeValueToBody(root, rootBody)
	}

	return f
}

// addScopeValueToBody writes a new block representing the scope value to the
// given body. This is especially useful for writing nested scope values.
func addScopeValueToBody(scope *nestedScope, body *hclwrite.Body) *hclwrite.Body {
	childBlock := body.AppendNewBlock("scope", []string{string(scope.Type), scope.Name})
	childBody := childBlock.Body()
	for _, grandchild := range scope.Children {
		childBody = addScopeValueToBody(grandchild, childBody)
	}
	return body
}

func (sg *scopeGenerator) writeScopeDataFile(b []byte) error {
	file, err := os.OpenFile(sg.filename, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			sg.Debug("file open errored, is ErrNotExist, creating file")
			file, err = os.Create(sg.filename)
			if err != nil {
				sg.Debug("file create failed")
				return err
			}
		} else {
			sg.Debug("file open errored, is not ErrNotExist, throwing")
			return err
		}
	} else {
		// err == nil means file was found
		sg.Warnf("A file '%s' already exists!", sg.filename)
		var yes bool
		survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("A file '%s' already exists! Overwrite?", sg.filename),
			Default: false,
		}, &yes)
		if !yes {
			sg.Infof("Not overwriting the existing file. Here is the generated scope data hcl:")
			fmt.Println(b)
			return nil
		}
	}
	_, err = file.Write(b)
	return err
}
