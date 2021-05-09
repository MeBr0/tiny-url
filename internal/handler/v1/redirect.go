package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/repo"
	"net/http"
)

func (h *Handler) initRedirectRoutes(api *gin.RouterGroup) {
	users := api.Group("/to")
	{
		users.GET("/:alias", h.redirectWithAlias)
	}
}

func (h *Handler) redirectWithAlias(c *gin.Context) {
	alias := c.Param("alias")

	if alias == "" {
		newResponse(c, http.StatusBadRequest, "empty alias")
		return
	}

	url, err := h.services.URLs.Get(c, alias)

	if err != nil {
		if err == repo.ErrURLNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.Original)
}
