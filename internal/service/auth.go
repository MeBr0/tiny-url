package service

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/auth"
	"github.com/mebr0/tiny-url/pkg/hash"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AuthService struct {
	repo           repo.Users
	hasher         hash.PasswordHasher
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func newAuthService(repo repo.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:           repo,
		hasher:         hasher,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, toRegister domain.UserRegister) error {
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

func (s *AuthService) Login(ctx context.Context, toLogin domain.UserLogin) (domain.Tokens, error) {
	passwordHash, err := s.hasher.Hash(toLogin.Password)

	if err != nil {
		return domain.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, toLogin.Email, passwordHash)

	if err != nil {
		return domain.Tokens{}, err
	}

	tokens, err := s.createSession(ctx, user.ID)

	// Async update last login
	if err == nil {
		go func() {
			c, cancel := context.WithTimeout(context.Background(), time.Duration(5) * time.Second)
			defer cancel()

			if err := s.repo.UpdateLastLogin(c, user.ID, time.Now()); err != nil {
				log.Warn("Could not updates last login by time.Now() for user " + user.ID.Hex() + " " + err.Error())
			}
		}()
	}

	return tokens, err
}

func (s *AuthService) createSession(ctx context.Context, userId primitive.ObjectID) (domain.Tokens, error) {
	var res domain.Tokens
	var err error

	res.AccessToken, err = s.tokenManager.Issue(userId.Hex(), s.accessTokenTTL)

	return res, err
}
