package service

import (
	"context"
	"github.com/mebr0/tiny-url/internal/cache"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/auth"
	"github.com/mebr0/tiny-url/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
}

type Auth interface {
	Register(ctx context.Context, toRegister domain.UserRegister) error
	Login(ctx context.Context, toLogin domain.UserLogin) (domain.Tokens, error)
}

type URLs interface {
	ListByOwner(ctx context.Context, owner primitive.ObjectID) ([]domain.URL, error)
	Create(ctx context.Context, toCreate domain.URLCreate) (domain.URL, error)
	Get(ctx context.Context, alias string) (domain.URL, error)
}

type Services struct {
	Users
	Auth
	URLs
}

func NewServices(repos *repo.Repos, caches *cache.Caches, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	urlEncoder hash.URLEncoder, accessTokenTTL time.Duration, aliasLength int, defaultExpiration int) *Services {
	return &Services{
		Users: newUsersService(repos.Users),
		Auth:  newAuthService(repos.Users, hasher, tokenManager, accessTokenTTL),
		URLs:  newURLsService(repos.URLs, caches.URLs, urlEncoder, aliasLength, defaultExpiration),
	}
}
