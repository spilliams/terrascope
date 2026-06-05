package pool

import (
	"bytes"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type Pool struct {
	workerCount int
	*logrus.Entry
}

func New(workerCount int, logger *logrus.Logger) *Pool {
	return &Pool{workerCount, logger.WithField("prefix", "pool")}
}

type TaskFunc func() ([]byte, error)

type taskOutput struct {
	err error
	log []byte
}

func (p *Pool) Execute(tasks []TaskFunc) ([]byte, error) {
	var wg sync.WaitGroup
	output := make(chan taskOutput, len(tasks))
	block := make(chan struct{}, p.workerCount)
	defer close(block)
	for _, task := range tasks {
		wg.Add(1)
		go func(task TaskFunc) {
			block <- struct{}{}
			p.Debugf("added task to pool")

			var out taskOutput
			out.log, out.err = task()
			output <- out

			p.Debugf("removing task from pool")
			<-block
			wg.Done()
		}(task)
	}
	wg.Wait()
	close(output)

	// compile the results
	var logs bytes.Buffer
	errs := make([]error, 0)
	for out := range output {
		_, err := logs.Write(out.log)
		if err != nil {
			errs = append(errs, err)
		}
		if out.err != nil {
			errs = append(errs, out.err)
		}
	}

	return logs.Bytes(), errors.Join(errs...)
}
