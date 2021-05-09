package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *Handler) initURLsRoutes(api *gin.RouterGroup) {
	users := api.Group("/urls", h.userIdentity)
	{
		users.GET("", h.listURLs)
		users.POST("", h.createURL)
	}
}

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

	urls, err := h.services.URLs.ListByOwner(c, userId)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	c.JSON(http.StatusOK, urls)
}

func (h *Handler) createURL(c *gin.Context) {
	var toCreate domain.URLCreate

	if err := c.BindJSON(&toCreate); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
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

	url, err := h.services.URLs.Create(c, toCreate)

	if err != nil {
		if err == repo.ErrURLAlreadyExists {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, url)
}
