package shell

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

func readStdoutAndStderr(log *logrus.Entry, stdout, stderr io.ReadCloser) error {
	out := newOutput()
	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	var stdoutErr, stderrErr error
	go func() {
		defer wg.Done()
		stdoutErr = readData(log, logrus.DebugLevel, stdoutReader, out.stdout)
	}()
	go func() {
		defer wg.Done()
		stderrErr = readData(log, logrus.WarnLevel, stderrReader, out.stderr)
	}()
	wg.Wait()

	if stdoutErr != nil {
		return stdoutErr
	}
	if stderrErr != nil {
		return stderrErr
	}

	return nil
}

func readData(log *logrus.Entry, logLevel logrus.Level, reader *bufio.Reader, writer io.StringWriter) error {
	var line string
	var readErr error
	for {
		line, readErr = reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		if len(line) == 0 && readErr == io.EOF {
			break
		}
		log.Log(logLevel, line)

		if _, err := writer.WriteString(line); err != nil {
			return err
		}
		if readErr != nil {
			break
		}
	}
	if readErr != io.EOF {
		return readErr
	}
	return nil
}

type output struct {
	stdout *stream
	stderr *stream
	merged *merged
}

func newOutput() *output {
	m := new(merged)
	return &output{
		merged: m,
		stdout: &stream{
			merged: m,
		},
		stderr: &stream{
			merged: m,
		},
	}
}

type merged struct {
	// no parallel writes
	sync.Mutex
	Lines []string
}

func (m *merged) WriteString(in string) (int, error) {
	m.Lock()
	defer m.Unlock()

	m.Lines = append(m.Lines, string(in))
	return len(in), nil
}

type stream struct {
	Lines []string
	*merged
}

func (s *stream) WriteString(in string) (int, error) {
	s.Lines = append(s.Lines, string(in))
	return s.merged.WriteString(in)
}
