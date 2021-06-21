package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	mockRepo "github.com/mebr0/tiny-url/internal/repo/mocks"
	"github.com/mebr0/tiny-url/pkg/auth"
	"github.com/mebr0/tiny-url/pkg/hash"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var commonErr = errors.New("error")

func mockAuthService(t *testing.T) (*AuthService, *mockRepo.MockUsers) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	usersRepo := mockRepo.NewMockUsers(mockCtl)
	authManager, _ := auth.NewManager("key")

	service := newAuthService(usersRepo, hash.NewSHA1Hasher(""), authManager, time.Duration(1) * time.Hour)

	return service, usersRepo
}

func TestAuthService_Register(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().Create(ctx, gomock.Any()).Return(primitive.NewObjectID(), nil)

	err := service.Register(ctx, domain.UserRegister{})

	require.NoError(t, err)
}

func TestAuthService_Login(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
	usersRepo.EXPECT().UpdateLastLogin(ctx, gomock.Any(), gomock.Any()).Return(nil)

	res, err := service.Login(ctx, domain.UserLogin{})

	require.NoError(t, err)
	require.IsType(t, domain.Tokens{}, res)
}

func TestAuthService_LoginErrUserNotExists(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, repo.ErrUserNotFound)

	_, err := service.Login(ctx, domain.UserLogin{})

	require.ErrorIs(t, err, repo.ErrUserNotFound)
}

func TestAuthService_LoginErr(t *testing.T) {
	service, usersRepo := mockAuthService(t)

	ctx := context.Background()

	usersRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any()).Return(domain.User{}, commonErr)

	_, err := service.Login(ctx, domain.UserLogin{})

	require.Error(t, err)
}
