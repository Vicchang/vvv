package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *Server) heartbeat(c *gin.Context) {
	resp := gin.H{
		"service":     "api",
		"status":      "ok",
		"connections": srv.connWatcher.Connections(),
	}

	c.JSON(http.StatusOK, resp)
}
