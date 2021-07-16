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
)

type URLsService struct {
	repo              repo.URLs
	cache             cache.URLs
	urlEncoder        hash.URLEncoder
	aliasLength       int
	defaultExpiration int
	urlCountLimit     int
}

func newURLsService(repo repo.URLs, cache cache.URLs, urlEncoder hash.URLEncoder, aliasLength int, defaultExpiration int, urlCountLimit int) *URLsService {
	return &URLsService{
		repo:              repo,
		cache:             cache,
		urlEncoder:        urlEncoder,
		aliasLength:       aliasLength,
		defaultExpiration: defaultExpiration,
		urlCountLimit:     urlCountLimit,
	}
}

func (s *URLsService) ListByOwner(ctx context.Context, userId primitive.ObjectID) ([]domain.URL, error) {
	return s.repo.ListByOwner(ctx, userId)
}

func (s *URLsService) ListByOwnerAndExpiration(ctx context.Context, userId primitive.ObjectID, expired bool) ([]domain.URL, error) {
	return s.repo.ListByOwnerAndExpiration(ctx, userId, expired)
}

func (s *URLsService) Create(ctx context.Context, toCreate domain.URLCreate) (domain.URL, error) {
	// Get URL from database
	_, err := s.repo.GetByOriginalAndOwner(ctx, toCreate.Original, toCreate.Owner)

	// If other error than not found return it
	if err != nil && err != repo.ErrURLNotFound {
		return domain.URL{}, err
	}

	// If URL exists return error
	if err == nil {
		return domain.URL{}, repo.ErrURLAlreadyExists
	}

	// Check for URL count limit
	urls, err := s.ListByOwner(ctx, toCreate.Owner)

	if err != nil {
		return domain.URL{}, err
	}

	if len(urls) > s.urlCountLimit {
		return domain.URL{}, ErrURLLimit
	}

	// Begin tries to generate alias
	try := 0

	for {
		alias, err := s.urlEncoder.Encode(toCreate.Original, toCreate.Owner, try, s.aliasLength)

		if err != nil {
			// Stop trying to create alias
			if err == hash.ErrURLAliasLengthExceed {
				break
			}

			return domain.URL{}, err
		}

		// Set default duration
		if toCreate.Duration == 0 {
			toCreate.Duration = s.defaultExpiration
		}

		url := domain.NewURL(toCreate, alias)
		id, err := s.repo.Create(ctx, url)

		if err != nil {
			if err != repo.ErrURLAlreadyExists {
				return domain.URL{}, err
			}

			// If alias already exists try one more
			log.Warn("Could not create alias")
			try += 1

			continue
		}

		return s.repo.Get(ctx, id)
	}

	return domain.URL{}, ErrNoPossibleAliasEncoding
}

func (s *URLsService) Get(ctx context.Context, alias string) (domain.URL, error) {
	// Get URL from cache
	url, err := s.cache.Get(ctx, alias)

	if err == nil {
		return url, nil
	}

	if err != redis.Nil {
		log.Warn("Error while get from cache " + err.Error())
	}

	// Get URL from database
	url, err = s.repo.Get(ctx, alias)

	if err != nil {
		return domain.URL{}, err
	}

	// Put valid URL to cache
	go func() {
		if err := s.cache.Set(ctx, url); err != nil {
			log.Warn("Could not save to cache " + err.Error())
		}
	}()

	return url, nil
}

func (s *URLsService) GetByOwner(ctx context.Context, alias string, owner primitive.ObjectID) (domain.URL, error) {
	url, err := s.Get(ctx, alias)

	if err != nil {
		return domain.URL{}, err
	}

	// If owners do not match, return forbidden
	if url.Owner != owner {
		return domain.URL{}, ErrURLForbidden
	}

	return url, nil
}

func (s *URLsService) Prolong(ctx context.Context, alias string, owner primitive.ObjectID, toProlong domain.URLProlong) (domain.URL, error) {
	if _, err := s.GetByOwner(ctx, alias, owner); err != nil {
		return domain.URL{}, err
	}

	if err := s.repo.Prolong(ctx, alias, toProlong); err != nil {
		return domain.URL{}, err
	}

	if err := s.cache.Delete(ctx, alias); err != nil {
		log.Warn("Error while delete from cache " + err.Error())
	}

	return s.GetByOwner(ctx, alias, owner)
}

func (s *URLsService) Delete(ctx context.Context, alias string, owner primitive.ObjectID) error {
	if _, err := s.GetByOwner(ctx, alias, owner); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, alias); err != nil {
		return err
	}

	if err := s.cache.Delete(ctx, alias); err != nil {
		log.Warn("Error while delete from cache " + err.Error())
	}

	return nil
}
