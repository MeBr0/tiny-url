package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *Handler) initURLsRoutes(api *gin.RouterGroup) {
	users := api.Group("/urls", h.userIdentity)
	{
		users.GET("", h.listURLs)
		users.POST("", h.createURL)
		users.GET("/:alias", h.getURL)
		users.PATCH("/:alias/prolong", h.prolongURL)
		users.DELETE("/:alias", h.deleteURL)
	}
}

// @Summary List URLs
// @Tags urls
// @Description List URLs owner by user
// @ID listURLs
// @Security UsersAuth
// @Accept json
// @Produce json
// @Success 200 {array} domain.URL "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 500 {object} response "Server error"
// @Router /urls [get]
func (h *Handler) listURLs(c *gin.Context) {
	userIdHex, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIdHex.(string))

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	urls, err := h.services.URLs.ListByOwner(c.Request.Context(), userId)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, urls)
}

// @Summary Create new URL
// @Tags urls
// @Description Create new URL for user
// @ID createURL
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param input body domain.URLCreate true "Data for creating URL"
// @Success 201 {object} domain.URL "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 422 {object} response "Invalid request body"
// @Failure 500 {object} response "Server error"
// @Router /urls [post]
func (h *Handler) createURL(c *gin.Context) {
	var toCreate domain.URLCreate

	if err := c.BindJSON(&toCreate); err != nil {
		newResponse(c, http.StatusUnprocessableEntity, "invalid request body")
		return
	}

	userIdHex, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIdHex.(string))

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	toCreate.Owner = userId

	url, err := h.services.URLs.Create(c.Request.Context(), toCreate)

	if err != nil {
		if err == repo.ErrURLAlreadyExists || err == service.ErrURLLimit {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, url)
}

// @Summary Get URL
// @Tags urls
// @Description Get URL by alias
// @ID getURL
// @Security UsersAuth
// @Accept json
// @Produce json
// @Success 200 {object} domain.URL "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /urls/{alias} [get]
func (h *Handler) getURL(c *gin.Context) {
	userIdHex, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIdHex.(string))

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	alias := c.Param("alias")

	if alias == "" {
		newResponse(c, http.StatusBadRequest, "empty alias")
		return
	}

	urls, err := h.services.URLs.GetByOwner(c.Request.Context(), alias, userId)

	if err != nil {
		if err == repo.ErrURLNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err == service.ErrURLForbidden {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, urls)
}

// @Summary Prolong URL
// @Tags urls
// @Description Prolong URL for user
// @ID prolongURL
// @Security UsersAuth
// @Accept json
// @Produce json
// @Param input body domain.URLProlong true "Data for prolonging URL"
// @Success 200 {object} domain.URL "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 422 {object} response "Invalid request body"
// @Failure 500 {object} response "Server error"
// @Router /urls/{alias}/prolong [patch]
func (h *Handler) prolongURL(c *gin.Context) {
	var toProlong domain.URLProlong

	if err := c.BindJSON(&toProlong); err != nil {
		newResponse(c, http.StatusUnprocessableEntity, "invalid request body")
		return
	}

	userIdHex, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIdHex.(string))

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	alias := c.Param("alias")

	if alias == "" {
		newResponse(c, http.StatusBadRequest, "empty alias")
		return
	}

	url, err := h.services.URLs.Prolong(c.Request.Context(), alias, userId, toProlong)

	if err != nil {
		if err == repo.ErrURLNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err == service.ErrURLForbidden {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, url)
}

// @Summary Delete URL
// @Tags urls
// @Description Delete URL by alias
// @ID deleteURL
// @Security UsersAuth
// @Accept json
// @Produce json
// @Success 204 {null} nil "Operation finished successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 401 {object} response "Invalid authorization"
// @Failure 403 {object} response "Invalid access"
// @Failure 500 {object} response "Server error"
// @Router /urls/{alias} [delete]
func (h *Handler) deleteURL(c *gin.Context) {
	userIdHex, ok := c.Get("userId")

	if !ok {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	userId, err := primitive.ObjectIDFromHex(userIdHex.(string))

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	alias := c.Param("alias")

	if alias == "" {
		newResponse(c, http.StatusBadRequest, "empty alias")
		return
	}

	if err := h.services.URLs.Delete(c.Request.Context(), alias, userId); err != nil {
		if err == repo.ErrURLNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if err == service.ErrURLForbidden {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
