package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/service"
)

var errUnsupportedContentType = errors.New("unsupported Content-Type header, only text/plain; charset=utf-8 allowed")
var maxBodySize = int64(2000)

func NewURLHandler() *urlHandler {
	return &urlHandler{service.NewURLService(), config.ServerConfig.Base.String()}
}

type urlHandler struct {
	service service.Service
	address string
}

func (h urlHandler) CreateURL(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		return errUnsupportedContentType
	}

	limitedBody := http.MaxBytesReader(w, r.Body, maxBodySize)
	urlBytes, err := io.ReadAll(limitedBody)
	defer limitedBody.Close()

	if err != nil {
		return fmt.Errorf("failed read request body: %w", err)
	}

	var id string
	id, err = h.service.CreateURL(urlBytes)

	if err != nil {
		return fmt.Errorf("failed create url from request body: %w", err)
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.address + "/" + id))

	if err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}

	return nil
}

func (h urlHandler) GetURL(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	url, err := h.service.GetURL(id)

	if err != nil {
		return fmt.Errorf("failed to get url by id: %w", err)
	}

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
