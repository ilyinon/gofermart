package services

import (
	"context"
	"errors"

	"gophermart/internal/entities"
	"gophermart/internal/repositories"
)

var ErrNotEnoughBalance = errors.New("not enough balance")

type BalanceService struct {
	orderRepo    repositories.OrderRepository
	withdrawRepo repositories.WithdrawalRepository
}

func NewBalanceService(
	orderRepo repositories.OrderRepository,
	withdrawRepo repositories.WithdrawalRepository,
) *BalanceService {

	return &BalanceService{
		orderRepo:    orderRepo,
		withdrawRepo: withdrawRepo,
	}
}

func (s *BalanceService) GetBalance(
	ctx context.Context,
	userID int64,
) (float64, float64, error) {

	orders, err := s.orderRepo.GetByUser(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	var accrual float64

	for _, o := range orders {
		accrual += o.Accrual
	}

	withdrawals, err := s.withdrawRepo.GetByUser(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	var withdrawn float64

	for _, w := range withdrawals {
		withdrawn += w.Sum
	}

	current := accrual - withdrawn

	return current, withdrawn, nil
}

func (s *BalanceService) Withdraw(
	ctx context.Context,
	userID int64,
	order string,
	sum float64,
) error {

	current, _, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	if current < sum {
		return ErrNotEnoughBalance
	}

	withdrawal := &entities.Withdrawal{
		UserID:      userID,
		OrderNumber: order,
		Sum:         sum,
	}

	return s.withdrawRepo.Create(ctx, withdrawal)
}

func (s *BalanceService) GetWithdrawals(
	ctx context.Context,
	userID int64,
) ([]entities.Withdrawal, error) {

	return s.withdrawRepo.GetByUser(ctx, userID)
}
