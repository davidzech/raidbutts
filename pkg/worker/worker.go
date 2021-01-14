package worker

import (
	"github.com/davidzech/raidbutts/pkg/job"
	"github.com/davidzech/raidbutts/pkg/simc"
)

type Options struct{}

// type Worker interface {
// 	Simulate(configuration simc.Configuration, options *Options) (simc.Result, error)
// }

type Worker struct {
	consumer job.Consumer
	maxJobs  int
	sim      simc.Simulator
}

type Option func(w *Worker)

func WithMaxJobs(n int) Option {
	return func(w *Worker) {
		w.maxJobs = n
	}
}

func WithSimulator(sim simc.Simulator) Option {
	return func(w *Worker) {
		w.sim = sim
	}
}

func NewWorker(consumer job.Consumer, opts ...Option) *Worker {
	w := Worker{
		consumer: consumer,
		maxJobs:  1,
		sim:      simc.DefautlCLI,
	}

	for _, opt := range opts {
		opt(&w)
	}

	return &w
}

type Result struct {
	Job    *job.Job
	Result *simc.Result
	Err    error
}

func (w *Worker) Start(resultsCh chan Result) error {
	jobsCh := make(chan *job.Job, w.maxJobs)

	go func() {
		j := <-jobsCh
		res, err := w.sim.Simulate(j.Config)
		resultsCh <- Result{
			Job:    j,
			Result: res,
			Err:    err,
		}
	}()

	for {
		job, err := w.consumer.Consume()
		if err != nil {
			return err
		}

		jobsCh <- job

		// TODO: terminate gracefully with select
	}
}

func (w *Worker) Stop() {
	panic("not implemented")
}
