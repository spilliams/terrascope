package terraboots

import (
	"fmt"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/internal/scopedata"
)

// NewScopeDataGenerator builds a generator for new scope data
func (p *Project) NewScopeDataGenerator(logger *logrus.Logger) (scopedata.Generator, error) {
	if len(p.Scopes) == 0 {
		return nil, fmt.Errorf("this project has no scope types! Please define them in %s with the terraboots `scope` block, then try this again", p.configFile)
	}

	scopeTypes := make([]string, len(p.Scopes))
	for i, el := range p.Scopes {
		scopeTypes[i] = el.Name
	}

	// this file doesn't have to exist yet
	dataFilename := "data.hcl"
	if p.ScopeData != nil && len(p.ScopeData) > 0 {
		// TODO: which filename? a new one? and then update the project config with the new filename?
		dataFilename = p.ScopeData[0]
	}
	scopeDataFile := path.Join(path.Dir(p.configFile), dataFilename)

	s := scopedata.NewGenerator(scopeTypes, scopeDataFile, logger)
	return s, nil
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
