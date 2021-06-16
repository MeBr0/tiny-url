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
	urls, err := s.repo.ListByOwner(ctx, userId)

	if err != nil {
		return nil, err
	}

	// Delete all expired URLs
	oneExpired := false

	for _, url := range urls {
		if url.Expired() {
			oneExpired = true

			if err := s.repo.Delete(ctx, url.Alias); err != nil {
				log.Warn("Could not delete from database " + err.Error())
			}
		}
	}

	// Query again if one expired
	if oneExpired {
		return s.repo.ListByOwner(ctx, userId)
	}

	return urls, nil
}

func (s *URLsService) Create(ctx context.Context, toCreate domain.URLCreate) (domain.URL, error) {
	// Get URL from database
	url, err := s.repo.GetByOriginalAndOwner(ctx, toCreate.Original, toCreate.Owner)

	// If other error than not found return it
	if err != nil && err != repo.ErrURLNotFound {
		return domain.URL{}, err
	}

	// If URL expired delete it
	if err == nil && url.Expired() {
		if err := s.repo.Delete(ctx, url.Alias); err != nil {
			log.Warn("Could not delete from database " + err.Error())
		}
	}

	// If URL exists return error
	if err == nil && !url.Expired() {
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

			log.Warn("Could not create alias")
			try += 1
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

	// If URL expired delete it from database and cache
	if url.Expired() {
		go func() {
			if err := s.cache.Delete(ctx, url.Alias); err != nil {
				log.Warn("Could not delete from cache " + err.Error())
			}
		}()

		go func() {
			if err := s.repo.Delete(ctx, url.Alias); err != nil {
				log.Warn("Could not delete from database " + err.Error())
			}
		}()

		return domain.URL{}, ErrURLExpired
	}

	// Put valid URL to cache
	go func() {
		if err := s.cache.Set(ctx, url); err != nil {
			log.Warn("Could not save to cache " + err.Error())
		}
	}()

	return url, nil
}
