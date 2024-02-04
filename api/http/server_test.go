package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEchoRoute(t *testing.T) {
	srv := NewServer("")
	router := srv.router

	expect := map[string]string{
		"echo": "echo",
		"test": "test",
	}

	data, err := json.Marshal(expect)
	require.NoError(t, err)

	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/echo", reader)
	router.ServeHTTP(w, req)

	actual := map[string]string{}
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expect, actual)
}

func TestHeartbeatRoute(t *testing.T) {
	srv := NewServer("")
	router := srv.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/heartbeat", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "api")
	assert.Contains(t, w.Body.String(), "ok")
}
