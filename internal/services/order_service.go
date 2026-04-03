package services

import (
	"context"
	"database/sql"
	"errors"

	"gophermart/internal/entities"
	"gophermart/internal/repositories"
)

var (
	ErrOrderExists = errors.New("order already uploaded by this user")
	ErrOrderUsed   = errors.New("order uploaded by another user")
)

type OrderService struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Upload(
	ctx context.Context,
	userID int64,
	number string,
) error {

	existing, err := s.repo.GetByNumber(ctx, number)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existing != nil {

		if existing.UserID == userID {
			return ErrOrderExists
		}

		return ErrOrderUsed
	}

	order := &entities.Order{
		Number: number,
		UserID: userID,
		Status: entities.StatusNew,
	}

	return s.repo.Create(ctx, order)
}

func (s *OrderService) List(
	ctx context.Context,
	userID int64,
) ([]entities.Order, error) {

	return s.repo.GetByUser(ctx, userID)
}
