package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/internal/service"
	"net/http"
)

func (h *Handler) initRedirectRoutes(api *gin.RouterGroup) {
	users := api.Group("/to")
	{
		users.GET("/:alias", h.redirectWithAlias)
	}
}

// @Summary Redirect
// @Tags urls
// @Description Redirect with alias
// @ID redirectWithAlias
// @Accept json
// @Produce json
// @Param alias path string true "Alias for redirection"
// @Success 301 {string} null "Redirected successfully"
// @Failure 400 {object} response "Invalid request"
// @Failure 500 {object} response "Server error"
// @Router /to/{alias} [get]
func (h *Handler) redirectWithAlias(c *gin.Context) {
	alias := c.Param("alias")

	if alias == "" {
		newResponse(c, http.StatusBadRequest, "empty alias")
		return
	}

	url, err := h.services.URLs.Get(c.Request.Context(), alias)

	if err != nil {
		if err == repo.ErrURLNotFound || err == service.ErrURLExpired {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.Original)
}
