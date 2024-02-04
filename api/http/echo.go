package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EchoEndpoint(c *gin.Context) {
	reader := c.Request.Body
	contentLength := c.Request.ContentLength
	contentType := c.Request.Header.Get("Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, nil)
}
