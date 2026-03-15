package repositories

import (
	"context"
	"gophermart/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByLogin(ctx context.Context, login string) (*entities.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	user *entities.User,
) error {

	query := `
	INSERT INTO users (login, password_hash)
	VALUES ($1, $2)
	RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Login,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	return err
}

func (r *userRepository) GetByLogin(
	ctx context.Context,
	login string,
) (*entities.User, error) {

	query := `
	SELECT id, login, password_hash, created_at
	FROM users
	WHERE login = $1
	`

	user := &entities.User{}

	err := r.db.QueryRow(
		ctx,
		query,
		login,
	).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
