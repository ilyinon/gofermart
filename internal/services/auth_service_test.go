package services

import (
	"context"
	"testing"

	"gophermart/internal/entities"
	"gophermart/internal/utils"
)

type mockUserRepo struct {
	user entities.User
}

func (m *mockUserRepo) Create(ctx context.Context, u *entities.User) error {
	m.user = *u
	m.user.ID = 1
	return nil
}

func (m *mockUserRepo) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	return &m.user, nil
}

func TestRegisterAndLogin(t *testing.T) {

	repo := &mockUserRepo{}

	jwt := utils.NewJWTManager("test-secret")

	service := NewAuthService(repo, jwt)

	err := service.Register(context.Background(), "test", "123")
	if err != nil {
		t.Fatal(err)
	}

	token, err := service.Login(context.Background(), "test", "123")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("token should not be empty")
	}
}
