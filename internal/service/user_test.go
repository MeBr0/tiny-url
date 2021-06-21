package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/mebr0/tiny-url/internal/domain"
	mockRepo "github.com/mebr0/tiny-url/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockUsersService(t *testing.T) (*UsersService, *mockRepo.MockUsers) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	usersRepo := mockRepo.NewMockUsers(mockCtl)

	service := newUsersService(usersRepo)

	return service, usersRepo
}

func TestUsersService_List(t *testing.T) {
	service, usersRepo := mockUsersService(t)

	ctx := context.Background()

	usersRepo.EXPECT().List(ctx).Return([]domain.User{}, nil)

	res, err := service.List(ctx)

	require.NoError(t, err)
	require.IsType(t, []domain.User{}, res)
}
