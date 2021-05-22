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

// @Summary Register
// @Tags auth
// @Description User registration
// @ID register
// @Accept json
// @Produce json
// @Param input body domain.UserRegister true "Register info"
// @Success 201 {string} null "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 422 {object} response "Invalid request body"
// @Failure 500 {object} response "Server error"
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var toRegister domain.UserRegister

	if err := c.BindJSON(&toRegister); err != nil {
		newResponse(c, http.StatusUnprocessableEntity, "invalid request body "+err.Error())
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

// @Summary Login
// @Tags auth
// @Description User login
// @ID login
// @Accept json
// @Produce json
// @Param input body domain.UserLogin true "Login credentials"
// @Success 200 {object} domain.Tokens "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 422 {object} response "Invalid request body"
// @Failure 500 {object} response "Server error"
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var toLogin domain.UserLogin

	if err := c.BindJSON(&toLogin); err != nil {
		newResponse(c, http.StatusUnprocessableEntity, "invalid request body "+err.Error())
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
