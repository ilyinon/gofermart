package services

import (
	"context"
	"testing"

	"gophermart/internal/entities"
)

type mockOrderRepoBalance struct{}

func (m *mockOrderRepoBalance) Create(ctx context.Context, o *entities.Order) error {
	return nil
}

func (m *mockOrderRepoBalance) GetByNumber(ctx context.Context, number string) (*entities.Order, error) {
	return nil, nil
}

func (m *mockOrderRepoBalance) GetByUser(ctx context.Context, userID int64) ([]entities.Order, error) {

	return []entities.Order{
		{Accrual: 100},
	}, nil
}

func (m *mockOrderRepoBalance) GetPending(ctx context.Context) ([]entities.Order, error) {
	return nil, nil
}

func (m *mockOrderRepoBalance) Update(ctx context.Context, o *entities.Order) error {
	return nil
}

type mockWithdrawalRepo struct{}

func (m *mockWithdrawalRepo) Create(ctx context.Context, w *entities.Withdrawal) error {
	return nil
}

func (m *mockWithdrawalRepo) GetByUser(ctx context.Context, userID int64) ([]entities.Withdrawal, error) {

	return []entities.Withdrawal{
		{Sum: 20},
	}, nil
}

func TestBalanceCalculation(t *testing.T) {

	orderRepo := &mockOrderRepoBalance{}
	withdrawRepo := &mockWithdrawalRepo{}

	service := NewBalanceService(orderRepo, withdrawRepo)

	current, withdrawn, err := service.GetBalance(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if current != 80 {
		t.Fatalf("expected current=80 got %f", current)
	}

	if withdrawn != 20 {
		t.Fatalf("expected withdrawn=20 got %f", withdrawn)
	}
}
