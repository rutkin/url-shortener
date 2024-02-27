package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/models"
	"github.com/rutkin/url-shortener/internal/app/service"
)

var errUnsupportedContentType = errors.New("unsupported Content-Type header, only text/plain; charset=utf-8 allowed")
var errUnsupportedBody = errors.New("unsupported body")
var maxBodySize = int64(2000)

func NewURLHandler() *urlHandler {
	return &urlHandler{service.NewURLService(), config.ServerConfig.Base.String()}
}

type urlHandler struct {
	service service.Service
	address string
}

func (h urlHandler) createResponseAddress(shortURL string) string {
	return h.address + "/" + shortURL
}

func (h urlHandler) CreateURL(w http.ResponseWriter, r *http.Request) error {
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
	_, err = w.Write([]byte(h.createResponseAddress(id)))

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

func (h urlHandler) CreateShorten(w http.ResponseWriter, r *http.Request) error {
	logger.Log.Info("create shorten")

	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	if len(req.URL) == 0 {
		return errUnsupportedBody
	}

	id, err := h.service.CreateURL([]byte(req.URL))

	if err != nil {
		return fmt.Errorf("failed create url from request body: %w", err)
	}

	resp := models.Response{
		Result: h.createResponseAddress(id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return fmt.Errorf("failed encode body: %w", err)
	}

	return nil
}
