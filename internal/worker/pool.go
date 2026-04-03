package worker

import (
	"context"
	"sync"
)

type Job func(ctx context.Context)

type Pool struct {
	jobs chan Job
	wg   sync.WaitGroup
}

func NewPool(size int, ctx context.Context) *Pool {

	p := &Pool{
		jobs: make(chan Job, 100),
	}

	for i := 0; i < size; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}

	return p
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			job(ctx)

		case <-ctx.Done():
			return
		}
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}
