package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/rutkin/url-shortener/internal/app/repository"
	"github.com/rutkin/url-shortener/internal/app/service"
	"github.com/stretchr/testify/assert"
)

var address = url.URL{Scheme: "http", Host: "localhost:8080"}

func Test_urlHandler_Main(t *testing.T) {
	tests := []struct {
		method       string
		path         string
		contentType  string
		requestBody  string
		expectedCode int
		expectedBody string
		location     string
	}{
		{method: http.MethodPost, path: "/unexpected", contentType: "text/plain; charset=utf-8", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodPost, path: "/", contentType: "unexpected", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodGet, path: "/notfound", contentType: "text/plain; charset=utf-8", expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodPost, path: "/", contentType: "text/plain; charset=utf-8", requestBody: "https://go.dev", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/D292748E"},
		{method: http.MethodGet, path: "/D292748E", contentType: "text/plain; charset=utf-8", expectedCode: http.StatusTemporaryRedirect, location: "https://go.dev"},
	}

	repository := repository.NewInMemoryRepository()

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			var reader io.Reader
			if tt.requestBody != "" {
				reader = strings.NewReader(tt.requestBody)
			}
			r := httptest.NewRequest(tt.method, tt.path, reader)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			service := service.NewURLService(repository)
			urlHandler := NewURLHandler(service, address)
			handler := MakeHandler(urlHandler.Main)

			handler(w, r)

			assert.Equal(t, tt.expectedCode, w.Code, "Response code does not match expected")
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String(), "Body is not equal to expected")
			}
			if tt.location != "" {
				assert.Equal(t, tt.location, w.HeaderMap.Get("Location"))
			}
		})
	}
}
