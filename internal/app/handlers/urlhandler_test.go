package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, contentType string) (*http.Response, string) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, ts.URL+path, reader)
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestURLHandlerRouter(t *testing.T) {
	ts := httptest.NewServer(NewURLHandlerRouter())
	defer ts.Close()

	tests := []struct {
		method       string
		path         string
		contentType  string
		requestBody  string
		expectedCode int
		expectedBody string
		location     string
	}{
		{method: http.MethodPost, path: "/unexpected", contentType: "text/plain; charset=utf-8", expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodPost, path: "/", contentType: "unexpected", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodGet, path: "/notfound", contentType: "text/plain; charset=utf-8", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodPost, path: "/", contentType: "text/plain; charset=utf-8", requestBody: "https://go.dev", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/D292748E"},
		{method: http.MethodGet, path: "/D292748E", contentType: "text/plain; charset=utf-8", expectedCode: http.StatusTemporaryRedirect, location: "https://go.dev"},
	}

	for _, tt := range tests {
		resp, get := testRequest(t, ts, tt.method, tt.path, tt.requestBody, tt.contentType)
		assert.Equal(t, tt.expectedCode, resp.StatusCode)
		if tt.expectedBody != "" {
			assert.Equal(t, tt.expectedBody, get)
		}
	}
}
