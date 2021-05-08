package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) ping(c *gin.Context) {
	log.Info(c.Get("userId"))

	c.JSON(http.StatusOK, nil)
}
