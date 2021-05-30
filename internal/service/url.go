package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mebr0/tiny-url/internal/cache"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/pkg/hash"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type URLsService struct {
	repo       repo.URLs
	cache      cache.URLs
	urlEncoder hash.URLEncoder
	aliasLength int
}

func newURLsService(repo repo.URLs, cache cache.URLs, urlEncoder hash.URLEncoder, aliasLength int) *URLsService {
	return &URLsService{
		repo:       repo,
		cache:      cache,
		urlEncoder: urlEncoder,
		aliasLength: aliasLength,
	}
}

func (s *URLsService) ListByOwner(ctx context.Context, userId primitive.ObjectID) ([]domain.URL, error) {
	return s.repo.ListByOwner(ctx, userId)
}

func (s *URLsService) Create(ctx context.Context, toCreate domain.URLCreate) (domain.URL, error) {
	try := 0

	for {
		alias, err := s.urlEncoder.Encode(toCreate.Original, toCreate.Owner, try, s.aliasLength)

		if err != nil {
			if err == hash.ErrURLAliasLengthExceed {
				break
			}

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
			if err != repo.ErrURLAlreadyExists {
				return domain.URL{}, err
			}

			log.Warn("Could not create alias")
			try += 1
		}

		return s.repo.Get(ctx, id)
	}

	return domain.URL{}, ErrNoPossibleAliasEncoding
}

func (s *URLsService) Get(ctx context.Context, alias string) (domain.URL, error) {
	url, err := s.cache.Get(ctx, alias)

	if err == nil {
		return url, nil
	}

	if err != redis.Nil {
		log.Warn("Error while get from cache " + err.Error())
	}

	url, err = s.repo.Get(ctx, alias)

	if err != nil {
		return domain.URL{}, err
	}

	go func() {
		if err := s.cache.Set(ctx, url); err != nil {
			log.Warn("Could not save to cache " + err.Error())
		}
	}()

	return url, nil
}
