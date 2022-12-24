package terraboots

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/internal/scopedata"
)

// GenerateScopeData builds a generator for new scope data, then executes it,
// saving the results in a file
func (p *Project) GenerateScopeData(input io.Writer, output io.Reader, logger *logrus.Logger) error {
	if len(p.Scopes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.Scopes))
	for i, el := range p.Scopes {
		scopeTypes[i] = el.Name
	}

	gen := scopedata.NewGenerator(scopeTypes, logger)
	err = gen.Create(input, output)
	if err != nil {
		return err
	}

	// this file doesn't have to exist yet
	dataFilename := "data.hcl"
	if p.ScopeData != nil && len(p.ScopeData) > 0 {
		// TODO: which filename? a new one? and then update the project config with the new filename?
		dataFilename = p.ScopeData[0]
	}
	scopeDataFile := path.Join(path.Dir(p.configFile), dataFilename)

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
	return err
}

// readScopeData reads all of the scope data known to the receiver
func (p *Project) readScopeData() error {
	if len(p.Scopes) == 0 {
		return fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	for _, filename := range p.ScopeData {
		err := p.readScopeDataFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// readScopeDataFile reads a single file with all of the receiver's scope data
// in it
func (p *Project) readScopeDataFile(filename string) error {
	scopeTypes := make([]string, len(p.Scopes))
	for i, el := range p.Scopes {
		scopeTypes[i] = el.Name
	}

	// WIP
	// gotta build some Specs?

	return nil
}
