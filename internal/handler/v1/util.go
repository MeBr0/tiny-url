package v1

import "github.com/gin-gonic/gin"

type response struct {
	// Success or error message
	Message string `json:"message" example:"some text"`
} // @name Response

func newResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{message})
}
