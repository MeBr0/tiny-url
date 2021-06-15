package app

import (
	"context"
	"errors"
	"github.com/mebr0/tiny-url/internal/cache"
	"github.com/mebr0/tiny-url/internal/config"
	"github.com/mebr0/tiny-url/internal/handler"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/internal/server"
	"github.com/mebr0/tiny-url/internal/service"
	"github.com/mebr0/tiny-url/pkg/auth"
	"github.com/mebr0/tiny-url/pkg/cache/redis"
	"github.com/mebr0/tiny-url/pkg/database/mongodb"
	"github.com/mebr0/tiny-url/pkg/hash"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Tiny URL API
// @version 1.0.0
// @description API for shortening URL

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes application
func Run(configPath string) {
	// Load configs
	cfg := config.LoadConfig(configPath)

	// Deps
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)

	if err != nil {
		log.Error(err)
		return
	}

	redisClient, err := redis.NewClient(cfg.Redis.URI, cfg.Redis.Password, cfg.Redis.Database)

	if err != nil {
		log.Error(err)
		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	passwordHasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
	urlHasher := hash.NewMD5Encoder()

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.Key)

	if err != nil {
		log.Error(err)
		return
	}

	// Init handlers
	repos := repo.NewRepos(db)
	caches := cache.NewCaches(redisClient, cfg.Redis.TTL)
	services := service.NewServices(repos, caches, passwordHasher, tokenManager, urlHasher, cfg.Auth.AccessTokenTTL, cfg.URL.AliasLength, cfg.URL.DefaultExpiration)
	handlers := handler.NewHandler(services, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
	}
}
