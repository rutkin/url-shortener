package app

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, contentType string) (int, string) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, ts.URL+path, reader)
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)

	ts.Client().CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(respBody)
}

func TestRootRouter(t *testing.T) {
	ts := httptest.NewServer(newRootRouter())
	defer ts.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		contentType  string
		requestBody  string
		expectedCode int
		expectedBody string
		location     string
	}{
		{
			name:         "method_post_unsupported_url",
			method:       http.MethodPost,
			path:         "/unexpected",
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_post_unsupported_content_type",
			method:       http.MethodPost,
			path:         "/",
			contentType:  "unexpected",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "method_get_unsupported_url",
			method:       http.MethodGet,
			path:         "/notfound",
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "method_post_success",
			method:       http.MethodPost,
			path:         "/",
			contentType:  "text/plain; charset=utf-8",
			requestBody:  "https://go.dev",
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080/D292748E",
		},
		{
			name:         "method_get_success",
			method:       http.MethodGet,
			path:         "/D292748E",
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusTemporaryRedirect,
			location:     "https://go.dev",
		},
		{
			name:         "method_post_shorten_unsupported_body",
			method:       http.MethodPost,
			path:         "/api/shorten",
			contentType:  "application/json",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "method_post_shorten_success",
			method:       http.MethodPost,
			path:         "/api/shorten",
			contentType:  "application/json",
			expectedCode: http.StatusCreated,
			requestBody:  `{"url": "https://testurl.com/blablabla"}`,
			expectedBody: `{"result":"http://localhost:8080/9718264F"}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := testRequest(t, ts, tt.method, tt.path, tt.requestBody, tt.contentType)
			assert.Equal(t, tt.expectedCode, status)

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestCompression(t *testing.T) {
	ts := httptest.NewServer(newRootRouter())
	defer ts.Close()

	requestBody := `{"url": "https://testurl.com/blablabla"}`
	expectedBody := `{"result":"http://localhost:8080/9718264F"}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", ts.URL+"/api/shorten", buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)
		require.JSONEq(t, expectedBody, string(b))
	})
}
