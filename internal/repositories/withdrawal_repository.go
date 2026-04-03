package repositories

import (
	"context"
	"fmt"

	"gophermart/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WithdrawalRepository interface {
	Create(
		ctx context.Context,
		withdrawal *entities.Withdrawal,
	) error

	GetByUser(
		ctx context.Context,
		userID int64,
	) ([]entities.Withdrawal, error)
}

type withdrawalRepository struct {
	*BaseRepository[entities.Withdrawal]
}

func NewWithdrawalRepository(db *pgxpool.Pool) WithdrawalRepository {
	return &withdrawalRepository{
		BaseRepository: NewBaseRepository[entities.Withdrawal](db),
	}
}

func (r *withdrawalRepository) Create(
	ctx context.Context,
	withdrawal *entities.Withdrawal,
) error {

	query := `
	INSERT INTO withdrawals (user_id, order_number, sum)
	VALUES ($1, $2, $3)
	RETURNING id, processed_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		withdrawal.UserID,
		withdrawal.OrderNumber,
		withdrawal.Sum,
	).Scan(
		&withdrawal.ID,
		&withdrawal.ProcessedAt,
	)

	if err != nil {
		return fmt.Errorf("create withdrawal: %w", err)
	}

	return nil
}

func (r *withdrawalRepository) GetByUser(
	ctx context.Context,
	userID int64,
) ([]entities.Withdrawal, error) {

	query := `
	SELECT id, user_id, order_number, sum, processed_at
	FROM withdrawals
	WHERE user_id = $1
	ORDER BY processed_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get withdrawals by user: %w", err)
	}

	defer rows.Close()

	var withdrawals []entities.Withdrawal

	for rows.Next() {

		var w entities.Withdrawal

		err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OrderNumber,
			&w.Sum,
			&w.ProcessedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan withdrawal: %w", err)
		}

		withdrawals = append(withdrawals, w)
	}

	return withdrawals, nil
}
