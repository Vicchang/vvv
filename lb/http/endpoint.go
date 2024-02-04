package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleEndpoint(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var data map[string]string
	err = json.Unmarshal(bs, &data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(data["url"])
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	fmt.Println(data["url"])
	s.PodService.Add(data["url"])
	c.Status(http.StatusOK)
}
