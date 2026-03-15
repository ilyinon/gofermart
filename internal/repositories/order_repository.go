package repositories

import (
	"context"
	"database/sql"

	"gophermart/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {

	Create(ctx context.Context, order *entities.Order) error

	GetByNumber(ctx context.Context, number string) (*entities.Order, error)

	GetByUser(ctx context.Context, userID int64) ([]entities.Order, error)

	GetPending(ctx context.Context) ([]entities.Order, error)

	Update(ctx context.Context, order *entities.Order) error
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(
	ctx context.Context,
	order *entities.Order,
) error {

	query := `
	INSERT INTO orders (number, user_id, status)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		order.Number,
		order.UserID,
		order.Status,
	)

	return err
}

func (r *orderRepository) GetByNumber(
	ctx context.Context,
	number string,
) (*entities.Order, error) {

	query := `
	SELECT number, user_id, status, accrual, uploaded_at
	FROM orders
	WHERE number = $1
	`

	order := &entities.Order{}

	var accrual sql.NullFloat64

	err := r.db.QueryRow(
		ctx,
		query,
		number,
	).Scan(
		&order.Number,
		&order.UserID,
		&order.Status,
		&accrual,
		&order.UploadedAt,
	)

	if err != nil {
		return nil, err
	}

	if accrual.Valid {
		order.Accrual = accrual.Float64
	}

	return order, nil
}

func (r *orderRepository) GetByUser(
	ctx context.Context,
	userID int64,
) ([]entities.Order, error) {

	query := `
	SELECT number, user_id, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []entities.Order

	for rows.Next() {

		var order entities.Order
		var accrual sql.NullFloat64

		err := rows.Scan(
			&order.Number,
			&order.UserID,
			&order.Status,
			&accrual,
			&order.UploadedAt,
		)

		if err != nil {
			return nil, err
		}

		if accrual.Valid {
			order.Accrual = accrual.Float64
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) GetPending(
	ctx context.Context,
) ([]entities.Order, error) {

	query := `
	SELECT number, user_id, status, accrual, uploaded_at
	FROM orders
	WHERE status IN ('NEW', 'PROCESSING')
	ORDER BY uploaded_at
	LIMIT 100
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []entities.Order

	for rows.Next() {

		var order entities.Order
		var accrual sql.NullFloat64

		err := rows.Scan(
			&order.Number,
			&order.UserID,
			&order.Status,
			&accrual,
			&order.UploadedAt,
		)

		if err != nil {
			return nil, err
		}

		if accrual.Valid {
			order.Accrual = accrual.Float64
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) Update(
	ctx context.Context,
	order *entities.Order,
) error {

	query := `
	UPDATE orders
	SET status = $2, accrual = $3
	WHERE number = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		order.Number,
		order.Status,
		order.Accrual,
	)

	return err
}
