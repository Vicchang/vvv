package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvv/lb/mock"
)

func TestEndpointRoute(t *testing.T) {
	gctrl := gomock.NewController(t)
	defer gctrl.Finish()
	mockPod := mock.NewMockPodService(gctrl)

	srv := NewServer()
	srv.PodService = mockPod

	router := srv.router

	expectURL := "http://localhost:8000"
	mockPod.EXPECT().Add(expectURL)

	data, err := json.Marshal(map[string]string{
		"url": expectURL,
	})
	require.NoError(t, err)

	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/endpoint", reader)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
