package service

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/auth"
	"github.com/mebr0/tiny-url/pkg/hash"
	"time"
)

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
}

type Auth interface {
	Register(ctx context.Context, toRegister domain.UserRegister) error
	Login(ctx context.Context, toLogin domain.UserLogin) (domain.Tokens, error)
}

type Services struct {
	Users
	Auth
}

func NewServices(repos *repo.Repos, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *Services {
	return &Services{
		Users: newUsersService(repos.Users),
		Auth: newAuthService(repos.Users, hasher, tokenManager, accessTokenTTL),
	}
}
