package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
