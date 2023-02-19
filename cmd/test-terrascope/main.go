package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/terrascope/internal/logformatter"
	"github.com/spilliams/terrascope/pkg/terrascope"
)

var log *logrus.Entry
var project *terrascope.Project

func init() {
	initLogger()
	if err := parseConfigAndProject(); err != nil {
		log.Panic(err)
	}
}

func initLogger() {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logformatter.PrefixedTextFormatter{
		UseColor: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	log = logger.WithField("prefix", "main")
}

func parseConfigAndProject() error {
	configFile := "../../../terrascope.hcl"
	log.Debugf("Using project configuration file: %s", configFile)
	var err error
	project, err = terrascope.ParseProject(configFile, log.Logger)
	if err != nil {
		return err
	}

	rootsDir := path.Join(path.Dir(configFile), project.RootsDir)
	rootsDir, err = filepath.Abs(rootsDir)
	if err != nil {
		return err
	}
	project.RootsDir = rootsDir
	// log.Debugf("Project roots directory: %s", project.RootsDir)
	// log.Debugf("Project scope data files: %s", project.ScopeDataFiles)

	return nil
}

func main() {
	rootName := "account-bucket"
	var scopes []string

	dirs, err := project.BuildRoot(rootName, scopes)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for _, dir := range dirs {
		log.Debugf(dir)
	}

}
