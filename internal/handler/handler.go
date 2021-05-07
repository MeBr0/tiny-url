package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/config"
	v1 "github.com/mebr0/tiny-url/internal/handler/v1"
	"github.com/mebr0/tiny-url/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	// Init router
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
