package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlPostRequest(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		fmt.Println("in valid method")
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	path := c.Request.RequestURI
	surl, err := s.SelectService.ServerURI()
	if err != nil {
		fmt.Println("no available server")
		c.Status(http.StatusInternalServerError)
		return
	}

	url, err := url.JoinPath(surl, path)
	if err != nil {
		fmt.Printf("in valid path, %v, %v\n", surl, path)
		c.Status(http.StatusInternalServerError)
		return
	}

	contentType := c.ContentType()
	body := c.Request.Body

	response, err := http.Post(url, contentType, body)
	if err != nil {
		fmt.Printf("post API error, %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		c.Status(response.StatusCode)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType = response.Header.Get("Content-Type")
	// TODO: fix manipulating of header values incorrectly
	extraHeader := map[string]string{}
	for k := range response.Header {
		extraHeader[k] = response.Header.Get(k)
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeader)
}
