package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvv/lb/mock"
)

func TestPostDispatchRoute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		w.Write(data)
	}))
	defer ts.Close()

	gctrl := gomock.NewController(t)
	defer gctrl.Finish()
	mockSelect := mock.NewMockSelectService(gctrl)
	mockSelect.EXPECT().ServerURI().Return(ts.URL, nil)

	srv := NewServer()
	srv.SelectService = mockSelect

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
