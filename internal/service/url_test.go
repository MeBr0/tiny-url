package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	mockCache "github.com/mebr0/tiny-url/internal/cache/mocks"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	mockRepo "github.com/mebr0/tiny-url/internal/repo/mocks"
	"github.com/mebr0/tiny-url/pkg/hash"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func mockURLService(t *testing.T) (*URLsService, *mockRepo.MockURLs, *mockCache.MockURLs) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	urlsRepo := mockRepo.NewMockURLs(mockCtl)
	urlsCache := mockCache.NewMockURLs(mockCtl)

	service := newURLsService(urlsRepo, urlsCache, hash.NewMD5Encoder(), 6, 10000, 3)

	return service, urlsRepo, urlsCache
}

func TestURLsService_ListByOwner(t *testing.T) {
	service, urlsRepo, _ := mockURLService(t)

	ctx := context.Background()

	userId := primitive.NewObjectID()

	urlsRepo.EXPECT().ListByOwner(ctx, userId).Return([]domain.URL{}, nil)

	res, err := service.ListByOwner(ctx, userId)

	require.NoError(t, err)
	require.IsType(t, []domain.URL{}, res)
}

func TestURLsService_ListByOwnerErr(t *testing.T) {
	service, urlsRepo, _ := mockURLService(t)

	ctx := context.Background()

	userId := primitive.NewObjectID()

	urlsRepo.EXPECT().ListByOwner(ctx, userId).Return([]domain.URL{}, commonErr)

	_, err := service.ListByOwner(ctx, userId)

	require.Error(t, err)
}

func TestURLsService_Create(t *testing.T) {
	service, urlsRepo, _ := mockURLService(t)

	ctx := context.Background()

	userId := primitive.NewObjectID()

	urlsRepo.EXPECT().GetByOriginalAndOwner(ctx, "url", userId).Return(domain.URL{}, repo.ErrURLNotFound)
	urlsRepo.EXPECT().ListByOwner(ctx, userId).Return([]domain.URL{}, nil)
	urlsRepo.EXPECT().Create(ctx, gomock.Any()).Return("alias", nil)
	urlsRepo.EXPECT().Get(ctx, "alias").Return(domain.URL{}, nil)

	res, err := service.Create(ctx, domain.URLCreate{
		Original: "url",
		Duration: 25,
		Owner:    userId,
	})

	require.NoError(t, err)
	require.IsType(t, domain.URL{}, res)
}

func TestURLsService_CreateErrURLAlreadyExists(t *testing.T) {
	service, urlsRepo, _ := mockURLService(t)

	ctx := context.Background()

	userId := primitive.NewObjectID()

	urlsRepo.EXPECT().GetByOriginalAndOwner(ctx, "url", userId).Return(domain.URL{
		ExpiredAt: time.Now().Add(time.Duration(1) * time.Minute),
	}, nil)

	_, err := service.Create(ctx, domain.URLCreate{
		Original: "url",
		Duration: 25,
		Owner:    userId,
	})

	require.ErrorIs(t, err, repo.ErrURLAlreadyExists)
}

func TestURLsService_GetFromCache(t *testing.T) {
	service, _, urlsCache := mockURLService(t)

	ctx := context.Background()

	urlsCache.EXPECT().Get(ctx, "alias").Return(domain.URL{}, nil)

	res, err := service.Get(ctx, "alias")

	require.NoError(t, err)
	require.IsType(t, domain.URL{}, res)
}

func TestURLsService_GetFromDatabase(t *testing.T) {
	service, urlsRepo, urlsCache := mockURLService(t)

	ctx := context.Background()

	urlsCache.EXPECT().Get(ctx, "alias").Return(domain.URL{}, redis.Nil)
	urlsRepo.EXPECT().Get(ctx, "alias").Return(domain.URL{
		ExpiredAt: time.Now().Add(time.Duration(1) * time.Minute),
	}, nil)
	urlsCache.EXPECT().Set(ctx, gomock.Any()).Return(nil)

	res, err := service.Get(ctx, "alias")

	require.NoError(t, err)
	require.IsType(t, domain.URL{}, res)
}
