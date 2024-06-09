package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func makeRequest(t *testing.B, ts *httptest.Server, method, path string, body string, contentType string) (int, string) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, ts.URL+path, reader)
	require.NoError(t, err)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept-Encoding", "identity")

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(respBody)
}

func BenchmarkCreateURLS(t *testing.B) {
	server, err := NewServer()
	require.NoError(t, err)
	defer server.Close()

	ts := httptest.NewServer(server.newRootRouter())
	defer ts.Close()

	for i := 0; i < 10000; i++ {
		var batch []string
		for j := 0; j < 100; j++ {
			batch = append(batch, fmt.Sprintf(`{"correlation_id": "%s","original_url": "https://testurl.com/%s"}`, uuid.NewString(), uuid.NewString()))
		}

		body := fmt.Sprintf("[%s]", strings.Join(batch, ","))

		status, resp := makeRequest(t, ts, http.MethodPost, "/api/shorten/batch", body, "application/json")
		fmt.Printf("status: %d", status)
		fmt.Printf("response: %s", resp)
	}
}
