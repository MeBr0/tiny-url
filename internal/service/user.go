package service

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/hash"
	"time"
)

type UsersService struct {
	repo repo.Users
	hasher hash.PasswordHasher
}

func newUsersService(repo repo.Users, hasher hash.PasswordHasher) *UsersService {
	return &UsersService{
		repo: repo,
		hasher: hasher,
	}
}

func (s *UsersService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UsersService) Create(ctx context.Context, toRegister domain.UserRegister) error {
	passwordHash, err := s.hasher.Hash(toRegister.Password)

	if err != nil {
		return err
	}

	user := domain.User{
		Name:         toRegister.Name,
		Email:        toRegister.Email,
		Password:     passwordHash,
		RegisteredAt: time.Now(),
		LastLogin:    time.Now(),
	}

	_, err = s.repo.Create(ctx, user)

	return err
}
