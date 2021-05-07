package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/domain"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.listUsers)
		users.POST("", h.createUser)
	}
}

func (h *Handler) listUsers(c *gin.Context) {
	users, err := h.services.Users.List(c)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) createUser(c *gin.Context) {
	var toRegister domain.UserRegister

	if err := c.BindJSON(&toRegister); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Create(c, toRegister); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
