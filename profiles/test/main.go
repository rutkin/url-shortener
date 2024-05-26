package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func makeRequest(url string, method, path string, body string, contentType string) (int, string) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url+path, reader)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept-Encoding", "identity")

	/*ts.Client().CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}*/
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return resp.StatusCode, string(respBody)
}

func main() {
	url := "http://localhost:8080"
	for i := 0; i < 10000; i++ {
		var batch []string
		for j := 0; j < 100; j++ {
			batch = append(batch, fmt.Sprintf(`{"correlation_id": "%s","original_url": "https://testurl.com/%s"}`, uuid.NewString(), uuid.NewString()))
		}

		body := fmt.Sprintf("[%s]", strings.Join(batch, ","))

		status, resp := makeRequest(url, http.MethodPost, "/api/shorten/batch", body, "application/json")
		fmt.Printf("status: %d", status)
		fmt.Printf("response: %s", resp)
	}
}
