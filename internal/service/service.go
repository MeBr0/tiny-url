package service

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/hash"
)

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, toRegister domain.UserRegister) error
}

type Services struct {
	Users
}

func NewServices(repos *repo.Repos, hasher hash.PasswordHasher) *Services {
	return &Services{
		Users: newUsersService(repos.Users, hasher),
	}
}
