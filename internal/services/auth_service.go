package services

import (
	"context"
	"errors"

	"gophermart/internal/entities"
	"gophermart/internal/repositories"
	"gophermart/internal/utils"
)

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

	return s.users.Create(ctx, user)
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
