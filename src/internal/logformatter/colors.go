package logformatter

import (
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

type colorScheme struct {
	TraceLevelColor string
	DebugLevelColor string
	InfoLevelColor  string
	WarnLevelColor  string
	ErrorLevelColor string
	FatalLevelColor string
	PanicLevelColor string
	PrefixColor     string
	TimestampColor  string
}

var (
	defaultColorScheme *colorScheme = &colorScheme{
		TraceLevelColor: "white",
		DebugLevelColor: "blue",
		InfoLevelColor:  "green",
		WarnLevelColor:  "yellow",
		ErrorLevelColor: "red",
		FatalLevelColor: "red",
		PanicLevelColor: "red",
		PrefixColor:     "cyan",
		TimestampColor:  "black+h",
	}
	noColorsColorScheme *colorScheme = &colorScheme{}
)

func (cs *colorScheme) colorFunc(s string) func(string) string {
	return ansi.ColorFunc(s)
}

func (cs *colorScheme) levelColorFunc(level logrus.Level) func(string) string {
	switch level {
	case logrus.TraceLevel:
		return cs.colorFunc(cs.TraceLevelColor)
	case logrus.DebugLevel:
		return cs.colorFunc(cs.DebugLevelColor)
	case logrus.InfoLevel:
		return cs.colorFunc(cs.InfoLevelColor)
	case logrus.WarnLevel:
		return cs.colorFunc(cs.WarnLevelColor)
	case logrus.ErrorLevel:
		return cs.colorFunc(cs.ErrorLevelColor)
	case logrus.FatalLevel:
		return cs.colorFunc(cs.FatalLevelColor)
	case logrus.PanicLevel:
		return cs.colorFunc(cs.PanicLevelColor)
	default:
		return cs.colorFunc(ansi.White)
	}
}

func (cs *colorScheme) prefixColorFunc() func(string) string {
	return cs.colorFunc(cs.PrefixColor)
}

func (cs *colorScheme) timestampColorFunc() func(string) string {
	return cs.colorFunc(cs.TimestampColor)
}
