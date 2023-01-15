package logformatter

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// PrefixedTextFormatter is the prefixed version of logrus.PrefixedTextFormatter
type PrefixedTextFormatter struct {
	UseColor bool

	PrintFullTimestamp bool
	// DisableTimestamp overrides PrintFullTimestamp
	DisableTimestamp bool
}

const miniTimestampFormat = "15:04:05.0000"

func (ptf *PrefixedTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	if entry.Buffer != nil {
		b = entry.Buffer
	}

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	timestampFormat := miniTimestampFormat
	if ptf.PrintFullTimestamp {
		timestampFormat = time.RFC3339
	}
	if ptf.DisableTimestamp {
		timestampFormat = ""
	}

	colorScheme := noColorsColorScheme
	if ptf.UseColor {
		colorScheme = defaultColorScheme
	}

	// this mutates b
	print(b, entry, keys, timestampFormat, colorScheme)

	return b.Bytes(), nil
}

func print(wr io.Writer, entry *logrus.Entry, keys []string,
	timestampFormat string, colorScheme *colorScheme) {

	levelColorFunc := colorScheme.levelColorFunc(entry.Level)
	levelText := entry.Level.String()
	if levelText == "warning" {
		levelText = "warn"
	}
	levelText = levelColorFunc(fmt.Sprintf(" %5s", levelText))

	prefixText := ""
	prefixColorFunc := colorScheme.prefixColorFunc()
	if prefixValue, ok := entry.Data["prefix"]; ok {
		prefixText = prefixColorFunc("[" + prefixValue.(string) + "] ")
	}

	timestampColorFunc := colorScheme.timestampColorFunc()
	timestampText := timestampColorFunc(entry.Time.Format(timestampFormat))

	fmt.Fprintf(wr, "%s%s %s%s", timestampText, levelText, prefixText, entry.Message)

	for _, k := range keys {
		if k == "prefix" {
			continue
		}

		v := entry.Data[k]
		fmt.Fprintf(wr, " %s=%+v", levelColorFunc(k), v)
	}

	fmt.Fprintln(wr, "")
}
