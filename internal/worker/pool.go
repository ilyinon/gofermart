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

func NewPool(size int) *Pool {

	p := &Pool{
		jobs: make(chan Job, 100),
	}

	for i := 0; i < size; i++ {

		p.wg.Add(1)

		go p.worker()
	}

	return p
}

func (p *Pool) worker() {

	defer p.wg.Done()

	for job := range p.jobs {

		job(context.Background())
	}
}

func (p *Pool) Submit(job Job) {

	p.jobs <- job
}

func (p *Pool) Stop() {

	close(p.jobs)

	p.wg.Wait()
}
