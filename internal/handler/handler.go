package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/config"
	v1 "github.com/mebr0/tiny-url/internal/handler/v1"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	// Init router
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler()

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
