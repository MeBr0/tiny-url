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

	service := newURLsService(urlsRepo, urlsCache, hash.NewMD5URLEncoder(), 6, 10000, 3)

	return service, urlsRepo, urlsCache
}

func TestURLsService_ListByOwnerAndExpiration(t *testing.T) {
	service, urlsRepo, _ := mockURLService(t)

	ctx := context.Background()

	userId := primitive.NewObjectID()

	urlsRepo.EXPECT().ListByOwnerAndExpiration(ctx, userId, false).Return([]domain.URL{}, nil)

	res, err := service.ListByOwnerAndExpiration(ctx, userId, false)

	require.NoError(t, err)
	require.IsType(t, []domain.URL{}, res)
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

	urlsRepo.EXPECT().ListByOwner(ctx, userId).Return([]domain.URL{}, errDefault)

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
	urlsCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)

	res, err := service.Get(ctx, "alias")

	require.NoError(t, err)
	require.IsType(t, domain.URL{}, res)
}

func TestURLsService_GetByOwner(t *testing.T) {
	s, _, urlsCache := mockURLService(t)

	ctx := context.Background()
	owner := primitive.NewObjectID()

	urlsCache.EXPECT().Get(ctx, "alias").Return(domain.URL{
		Owner: owner,
	}, nil)

	res, err := s.GetByOwner(ctx, "alias", owner)

	require.NoError(t, err)
	require.IsType(t, domain.URL{}, res)
}

func TestURLsService_GetByOwnerErrURLForbidden(t *testing.T) {
	s, _, urlsCache := mockURLService(t)

	ctx := context.Background()

	urlsCache.EXPECT().Get(ctx, "alias").Return(domain.URL{
		Owner: primitive.NilObjectID,
	}, nil)

	_, err := s.GetByOwner(ctx, "alias", primitive.NewObjectID())

	require.ErrorIs(t, err, ErrURLForbidden)
}

func TestURLsService_Prolong(t *testing.T) {
	s, urlsRepo, urlsCache := mockURLService(t)

	ctx := context.Background()

	owner := primitive.NewObjectID()

	urlsCache.EXPECT().Get(ctx, "alias").Return(domain.URL{
		Owner: owner,
	}, nil).Times(2)

	urlsRepo.EXPECT().Prolong(ctx, "alias", domain.URLProlong{Duration: 5}).Return(nil)
	urlsCache.EXPECT().Delete(gomock.Any(), "alias").Return(nil)

	res, err := s.Prolong(ctx, "alias", owner, domain.URLProlong{Duration: 5})

	require.NoError(t, err)
	require.IsType(t, domain.URL{}, res)
}

func TestURLsService_Delete(t *testing.T) {
	s, urlsRepo, urlsCache := mockURLService(t)

	ctx := context.Background()

	owner := primitive.NewObjectID()

	urlsCache.EXPECT().Get(ctx, "alias").Return(domain.URL{
		Owner: owner,
	}, nil)

	urlsRepo.EXPECT().Delete(ctx, "alias").Return(nil)
	urlsCache.EXPECT().Delete(gomock.Any(), "alias").Return(nil)

	err := s.Delete(ctx, "alias", owner)

	require.NoError(t, err)
}
