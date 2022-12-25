package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spilliams/terraboots/pkg/logformatter"
)

var log = logrus.New()

func init() {
	log.Formatter = &logformatter.PrefixedTextFormatter{
		UseColor: true,
	}
	log.Level = logrus.TraceLevel
}

func main() {
	log.Trace("this is a trace message")

	log.WithFields(logrus.Fields{
		"prefix": "main",
		"foo":    "bar",
	}).Debug("this is a debug message")

	log.WithFields(logrus.Fields{
		"prefix": "sensor",
	}).Info("this is an info message")

	log.Warn("this is a warn message")

	log.Error("This is an error message")

	log.Print("this is a print message")
}
