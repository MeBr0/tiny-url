package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"net/http"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	users := api.Group("/auth")
	{
		users.POST("/register", h.register)
		users.POST("/login", h.login)
	}
}

func (h *Handler) register(c *gin.Context) {
	var toRegister domain.UserRegister

	if err := c.BindJSON(&toRegister); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Register(c, toRegister); err != nil {
		if err == repo.ErrUserAlreadyExists {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) login(c *gin.Context) {
	var toLogin domain.UserLogin

	if err := c.BindJSON(&toLogin); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.services.Login(c, toLogin)

	if err != nil {
		if err == repo.ErrUserNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONP(http.StatusOK, token)
}
