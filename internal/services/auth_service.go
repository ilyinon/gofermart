package services

import (
	"context"
	"errors"

	"gophermart/internal/entities"
	"gophermart/internal/repositories"
	"gophermart/internal/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserExists = errors.New("user exists")

type AuthService struct {
	users repositories.UserRepository
}

func NewAuthService(users repositories.UserRepository) *AuthService {
	return &AuthService{
		users: users,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	login string,
	password string,
) error {

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &entities.User{
		Login:    login,
		Password: hash,
	}

	err = s.users.Create(ctx, user)
	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ErrUserExists
			}
		}

		return err
	}

	return nil
}

func (s *AuthService) Login(
	ctx context.Context,
	login string,
	password string,
) (string, error) {

	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
