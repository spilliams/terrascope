package pool

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"

	"github.com/sirupsen/logrus"
	"github.com/spilliams/terrascope/internal/logformatter"
)

func logger() *logrus.Logger {
	l := logrus.StandardLogger()
	l.SetFormatter(&logformatter.PrefixedTextFormatter{})
	l.SetLevel(logrus.DebugLevel)
	return l
}

func TestExecuteInSerial(t *testing.T) {
	words := []string{"foo", "bar", "baz", "qux"}
	testExecute(t, words, 1)
}

func TestExecuteInParallel(t *testing.T) {
	words := []string{"foo", "bar", "baz", "qux"}
	testExecute(t, words, 4)
}

func TestExecuteOverlargePool(t *testing.T) {
	words := []string{"foo", "bar", "baz", "qux"}
	testExecute(t, words, 400)
}

func testExecute(t *testing.T, words []string, workerCount int) {
	p := New(workerCount, logger())

	tasks := make([]TaskFunc, len(words))
	numLogLines := 0
	for i, word := range words {
		var n int
		n, tasks[i] = makeTask(word, false)
		numLogLines += n
	}

	logs, err := p.Execute(tasks)
	lines := strings.Split(string(logs), "\n")
	assert.Assert(t, is.Len(lines, numLogLines+1))
	assert.NilError(t, err)
}

func TestExecuteErrors(t *testing.T) {
	p := New(4, logger())

	words := []string{"foo", "bar", "baz", "qux"}
	tasks := make([]TaskFunc, len(words))
	numLogLines := 0
	for i, word := range words {
		var n int
		n, tasks[i] = makeTask(word, true)
		numLogLines += n
	}

	_, err := p.Execute(tasks)
	assert.ErrorContains(t, err, "uh oh")
	// assert.Assert(t, is.Len(errs, len(tasks)))
}

// returns a TaskFunc and the number of lines it logs
func makeTask(word string, makeError bool) (int, TaskFunc) {
	return 2, func() ([]byte, error) {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("beginning of task '%s'\n", word))
		time.Sleep(100 * time.Millisecond)
		buf.WriteString(fmt.Sprintf("end of task '%s'\n", word))
		var err error
		if makeError {
			err = fmt.Errorf("uh oh, '%s' had an error", word)
		}
		return buf.Bytes(), err
	}
}
