package service

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type URLsService struct {
	repo       repo.URLs
	urlEncoder hash.URLEncoder
}

func newURLsService(repo repo.URLs, urlEncoder hash.URLEncoder) *URLsService {
	return &URLsService{
		repo:       repo,
		urlEncoder: urlEncoder,
	}
}

func (s *URLsService) ListByOwner(ctx context.Context, userId primitive.ObjectID) ([]domain.URL, error) {
	return s.repo.ListByOwner(ctx, userId)
}

func (s *URLsService) Create(ctx context.Context, toCreate domain.URLCreate) (domain.URL, error) {
	alias, err := s.urlEncoder.Encode(toCreate.Original)

	if err != nil {
		return domain.URL{}, err
	}

	url := domain.URL{
		Alias:     alias,
		Original:  toCreate.Original,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(10000000000)),
		Owner:     toCreate.Owner,
	}

	id, err := s.repo.Create(ctx, url)

	if err != nil {
		return domain.URL{}, err
	}

	return s.repo.Get(ctx, id)
}

func (s *URLsService) Get(ctx context.Context, alias string) (domain.URL, error) {
	return s.repo.Get(ctx, alias)
}
