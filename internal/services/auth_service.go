package services

import (
	"context"
	"errors"
	"fmt"

	"gophermart/internal/entities"
	"gophermart/internal/repositories"
	"gophermart/internal/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserExists         = errors.New("user exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	users repositories.UserRepository
	jwt   *utils.JWTManager
}

func NewAuthService(
	users repositories.UserRepository,
	jwt *utils.JWTManager,
) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwt,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	login string,
	password string,
) error {

	hash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
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

		return fmt.Errorf("create user: %w", err)
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
		return "", fmt.Errorf("get user: %w", err)
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
