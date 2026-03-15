package services

import (
	"context"
	"testing"

	"gophermart/internal/entities"
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

	service := NewAuthService(repo)

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
