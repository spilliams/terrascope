package terraboots

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func (p *Project) NewScopeDataGenerator(logger *logrus.Logger) (ScopeDataGenerator, error) {
	if len(p.Scopes) == 0 {
		return nil, fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.Scopes))
	for i, el := range p.Scopes {
		scopeTypes[i] = el.Name
	}

	// this file doesn't have to exist yet
	scopeDataFile := path.Join(p.configFile, p.ScopeData)

	s := &scopeDataGenerator{
		scopeTypes: scopeTypes,
		fileName:   scopeDataFile,
		// TODO: read existing data to put in here
		// scopeData:

		Logger: logger,
	}

	return s, nil
}

type ScopeDataGenerator interface {
	Create(io.Reader, io.Writer) error
}

type scopeDataGenerator struct {
	scopeTypes []string
	scopeData  map[string]interface{}
	fileName   string
	*logrus.Logger
}

type scopeValue struct {
	address        string
	scopeTypeIndex int
}

func (sdg *scopeDataGenerator) Create(input io.Reader, output io.Writer) error {
	scopes, err := sdg.promptForScopeValues(input, output)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		sdg.Warn("No scopes were generated, exiting.")
		return nil
	}

	for _, scope := range scopes {
		sdg.Debugf("scope address: %s", scope)
	}
	sdg.Infof("%d scope addresses created", len(scopes))

	file, err := sdg.generateScopeDataFile(scopes)
	if err != nil {
		return err
	}

	sdg.Debug(string(file))
	// TODO write the file

	sdg.Warn("the rest of this is not yet implemented")

	return nil
}

func (sdg *scopeDataGenerator) promptForScopeValues(input io.Reader, output io.Writer) ([]string, error) {
	logrus.Debugf("[scopeDataGenerator.Create]")

	fmt.Fprintln(output, "Scope types in this projct, in order, are:")
	fmt.Fprintln(output, strings.Join(sdg.scopeTypes, ", "))
	fmt.Fprintln(output, "")

	fmt.Fprintln(output, "Answers must be space-separated, and may consist of the characters")
	// TODO: use this charset to validate the input...
	answerCharacterSet := "0-9a-zA-Z-_"
	fmt.Fprintln(output, answerCharacterSet)
	fmt.Fprintln(output, "")

	fmt.Fprintln(output, "Leave any answer blank to mark the current scope as complete with no children")
	fmt.Fprintln(output, "")

	fmt.Fprintln(output, "Press Ctrl+C at any time to cancel.")
	fmt.Fprintln(output, "")

	scanner := bufio.NewScanner(input)

	// First one's free
	fmt.Fprintf(output, "What are the allowable values for `%s`?\n", sdg.scopeTypes[0])
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	if len(scanner.Text()) == 0 {
		sdg.Debugf("user entered blank line, exiting")
		return nil, nil
	}
	// TODO: validate input against list of blocklisted words, and the above
	// charset, and each other (no dupes)...
	firstValues := strings.Split(scanner.Text(), " ")
	sdg.Debugf("read new scope values %v", firstValues)

	scopes := make([]string, 0)

	prompts := make([]scopeValue, len(firstValues))
	for i, el := range firstValues {
		prompts[i] = scopeValue{
			address:        el,
			scopeTypeIndex: 0,
		}
	}

	for len(prompts) > 0 {
		prompt := prompts[0]
		prompts = prompts[1:]

		if prompt.scopeTypeIndex+1 == len(sdg.scopeTypes) {
			// this scope value cannot have children
			scopes = append(scopes, prompt.address)
			continue
		}

		fmt.Fprintf(output, "Within %s, what are the allowable values for `%s`?\n", prompt.address, sdg.scopeTypes[prompt.scopeTypeIndex+1])

		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			return nil, err
		}
		if len(scanner.Text()) == 0 {
			// user entered none
			sdg.Debugf("user entered no values for prompt, closing this scope")
			scopes = append(scopes, prompt.address)
			continue
		}
		// TODO: validate input against list of blocklisted words, and the above
		// charset, and each other (no dupes)...

		values := strings.Split(scanner.Text(), " ")
		sdg.Debugf("read new scope values %v", values)
		for _, el := range values {
			value := scopeValue{
				address:        strings.Join([]string{prompt.address, el}, "."),
				scopeTypeIndex: prompt.scopeTypeIndex + 1,
			}
			prompts = append(prompts, value)
		}
	}

	return scopes, nil
}

func (sdg *scopeDataGenerator) generateScopeDataFile(scopes []string) ([]byte, error) {
	// TODO: check out hclwrite!
	// https://pkg.go.dev/github.com/hashicorp/hcl/v2@v2.15.0/hclwrite

	return nil, nil
}
