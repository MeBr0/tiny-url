package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/service"
	"github.com/mebr0/tiny-url/pkg/auth"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initAuthRoutes(v1)
		h.initURLsRoutes(v1)
		h.initRedirectRoutes(v1)

		v1.GET("/ping", h.userIdentity, h.ping)
	}
}
