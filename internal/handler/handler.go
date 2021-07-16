package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/config"
	v1 "github.com/mebr0/tiny-url/internal/handler/v1"
	"github.com/mebr0/tiny-url/internal/service"
	"github.com/mebr0/tiny-url/pkg/auth"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	// Init swagger routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init router
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
