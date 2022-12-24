package scopedata

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
)

type Generator interface {
	Create(io.Reader, io.Writer) error
}

func NewGenerator(scopeTypes []string, filename string, logger *logrus.Logger) Generator {
	return &generator{
		scopeTypes: scopeTypes,
		filename:   filename,
		Logger:     logger,
	}
}

type generator struct {
	scopeTypes []string
	filename   string
	*logrus.Logger
}

type scopeValue struct {
	name           string
	scopeType      string
	address        string
	children       map[string]scopeValue
	scopeTypeIndex int
}

func (g *generator) Create(input io.Reader, output io.Writer) error {
	rootScopes, err := g.promptForScopeValues(input, output)
	if err != nil {
		return err
	}
	if len(rootScopes) == 0 {
		g.Warn("No scopes were generated, exiting.")
		return nil
	}

	hclfile := g.generateScopeDataFile(rootScopes)

	file, err := os.OpenFile(g.filename, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			g.Debug("file open errored, is ErrNotExist, creating file")
			file, err = os.Create(g.filename)
			if err != nil {
				g.Debug("file create failed")
				return err
			}
		} else {
			g.Debug("file open errored, is not ErrNotExist, throwing")
			return err
		}
	} else {
		// err == nil means file was found
		g.Warnf("A file '%s' already exists! Overwrite? [Y/n]", g.filename)
		scanner := bufio.NewScanner(input)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			g.Debug("scanner errored")
			return err
		}
		if len(scanner.Text()) != 0 {
			g.Debug("scanner returned text")
			if scanner.Text() != "y" && scanner.Text() != "Y" {
				g.Debugf("User does not want to overwrite, printing and exiting.")
				output.Write(hclfile.Bytes())
				return nil
			}
		}
	}
	_, err = hclfile.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

func (g *generator) promptForScopeValues(input io.Reader, output io.Writer) ([]scopeValue, error) {
	logrus.Debugf("[scopeDataGenerator.Create]")

	fmt.Fprintln(output, "Scope types in this projct, in order, are:")
	fmt.Fprintln(output, strings.Join(g.scopeTypes, ", "))
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
	fmt.Fprintf(output, "What are the allowable values for `%s`?\n", g.scopeTypes[0])
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	if len(scanner.Text()) == 0 {
		g.Debugf("user entered blank line, exiting")
		return nil, nil
	}
	// TODO: validate input against list of blocklisted words, and the above
	// charset, and each other (no dupes)...
	firstValues := strings.Split(scanner.Text(), " ")
	g.Debugf("read new scope values %v", firstValues)

	roots := make([]scopeValue, len(firstValues))
	prompts := make([]scopeValue, len(firstValues))
	for i, el := range firstValues {
		value := scopeValue{
			name:           el,
			scopeType:      g.scopeTypes[0],
			scopeTypeIndex: 0,
			children:       make(map[string]scopeValue),
		}
		value.address = fmt.Sprintf("%s.%s", value.scopeType, value.name)
		roots[i] = value
		prompts[i] = value
	}

	for len(prompts) > 0 {
		prompt := prompts[0]
		prompts = prompts[1:]

		if prompt.scopeTypeIndex+1 == len(g.scopeTypes) {
			// this scope value cannot have children
			continue
		}

		fmt.Fprintf(output, "Within %s, what are the allowable values for `%s`?\n", prompt.address, g.scopeTypes[prompt.scopeTypeIndex+1])

		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			return nil, err
		}
		if len(scanner.Text()) == 0 {
			// user entered none
			g.Debugf("user entered no values for prompt, closing this scope")
			continue
		}
		// TODO: validate input against list of blocklisted words, and the above
		// charset, and each other (no dupes)...

		values := strings.Split(scanner.Text(), " ")
		g.Debugf("read new scope values %v", values)
		for _, el := range values {
			value := scopeValue{
				name:           el,
				scopeType:      g.scopeTypes[prompt.scopeTypeIndex+1],
				scopeTypeIndex: prompt.scopeTypeIndex + 1,
				children:       make(map[string]scopeValue),
			}
			value.address = strings.Join([]string{prompt.address, value.scopeType, value.name}, ".")
			prompt.children[el] = value
			prompts = append(prompts, value)
		}
	}

	g.Debugf("%+v", roots)

	return roots, nil
}

func (g *generator) generateScopeDataFile(rootScopes []scopeValue) *hclwrite.File {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	rootBody.AppendUnstructuredTokens(commentTokens("This file was generated by terraboots"))
	rootBody.AppendNewline()

	for _, root := range rootScopes {
		rootBody = addScopeValueToBody(root, rootBody)
	}

	return f
}

func addScopeValueToBody(scope scopeValue, body *hclwrite.Body) *hclwrite.Body {
	childBlock := body.AppendNewBlock(scope.scopeType, []string{scope.name})
	childBody := childBlock.Body()
	for _, grandchild := range scope.children {
		childBody = addScopeValueToBody(grandchild, childBody)
	}
	return body
}

func commentTokens(msg string) hclwrite.Tokens {
	if !strings.HasPrefix(msg, "# ") {
		msg = fmt.Sprintf("# %s", msg)
	}
	msgToken := &hclwrite.Token{
		Type:         hclsyntax.TokenComment,
		Bytes:        []byte(msg),
		SpacesBefore: 0,
	}
	return []*hclwrite.Token{msgToken}
}
