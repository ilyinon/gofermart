package worker

import (
	"context"
	"time"

	"gophermart/internal/entities"
	"gophermart/internal/infrastructure"
	"gophermart/internal/repositories"
)

type AccrualWorker struct {
	repo    repositories.OrderRepository
	accrual infrastructure.AccrualClient
	pool    *Pool

	stop chan struct{}
}

func NewAccrualWorker(
	repo repositories.OrderRepository,
	accrual infrastructure.AccrualClient,
) *AccrualWorker {

	return &AccrualWorker{
		repo:    repo,
		accrual: accrual,
		pool:    NewPool(5),
		stop:    make(chan struct{}),
	}
}

func (w *AccrualWorker) Start(ctx context.Context) {

	ticker := time.NewTicker(5 * time.Second)

	for {

		select {

		case <-ticker.C:

			orders, err := w.repo.GetPending(ctx)
			if err != nil {
				continue
			}

			for _, order := range orders {

				o := order

				w.pool.Submit(func(ctx context.Context) {

					resp, err := w.accrual.GetOrder(ctx, o.Number)

					if err != nil || resp == nil {
						return
					}

					if resp.Status == entities.StatusProcessed ||
						resp.Status == entities.StatusInvalid {

						o.Status = resp.Status
						o.Accrual = resp.Accrual

						w.repo.Update(ctx, &o)
						return
					}

					if resp.Status == entities.StatusProcessing {

						o.Status = resp.Status

						w.repo.Update(ctx, &o)
					}
				})
			}

		case <-w.stop:

			ticker.Stop()

			w.pool.Stop()

			return
		}
	}
}

func (w *AccrualWorker) Stop() {

	close(w.stop)
}
