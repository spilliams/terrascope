package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/internal/terraboots"
	"github.com/spilliams/terraboots/pkg/logformatter"
)

var log *logrus.Entry
var project *terraboots.Project

func init() {
	initLogger()
	if err := bootsbootsPreRunE(); err != nil {
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

func bootsbootsPreRunE() error {
	configFile := "../../../terraboots.hcl"
	log.Debugf("Using project configuration file: %s", configFile)
	var err error
	project, err = terraboots.ParseProject(configFile, log.Logger)
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
