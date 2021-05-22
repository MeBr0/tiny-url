package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Ping
// @Tags utils
// @Description Simple ping
// @ID ping
// @Accept json
// @Produce json
// @Success 200 {string} null "Operation finished successfully"
// @Router /ping [get]
func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
