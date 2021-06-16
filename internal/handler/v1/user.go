package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.listUsers)
	}
}

// @Summary List users
// @Tags users
// @Description List all users
// @ID listUsers
// @Accept json
// @Produce json
// @Success 200 {array} domain.User "Operation finished successfully"
// @Failure 500 {object} response "Server error"
// @Router /users [get]
func (h *Handler) listUsers(c *gin.Context) {
	users, err := h.services.Users.List(c.Request.Context())

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}
