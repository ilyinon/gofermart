package services

import (
	"context"
	"testing"

	"gophermart/internal/entities"
)

type mockOrderRepo struct {
	orders []entities.Order
}

func (m *mockOrderRepo) Create(ctx context.Context, o *entities.Order) error {
	m.orders = append(m.orders, *o)
	return nil
}

func (m *mockOrderRepo) GetByNumber(ctx context.Context, number string) (*entities.Order, error) {

	for _, o := range m.orders {
		if o.Number == number {
			return &o, nil
		}
	}

	return nil, nil
}

func (m *mockOrderRepo) GetByUser(ctx context.Context, userID int64) ([]entities.Order, error) {
	return m.orders, nil
}

func (m *mockOrderRepo) GetPending(ctx context.Context) ([]entities.Order, error) {
	return nil, nil
}

func (m *mockOrderRepo) Update(ctx context.Context, o *entities.Order) error {
	return nil
}

func TestUploadOrder(t *testing.T) {

	repo := &mockOrderRepo{}

	service := NewOrderService(repo)

	err := service.Upload(context.Background(), 1, "79927398713")
	if err != nil {
		t.Fatal(err)
	}

	if len(repo.orders) != 1 {
		t.Fatal("order not created")
	}
}
