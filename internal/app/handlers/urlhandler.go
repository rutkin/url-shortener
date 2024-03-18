package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/models"
	"github.com/rutkin/url-shortener/internal/app/service"
	"go.uber.org/zap"
)

var errUnsupportedBody = errors.New("unsupported body")
var maxBodySize = int64(2000)

func NewURLHandler() (*URLHandler, error) {
	s, err := service.NewURLService()
	if err != nil {
		logger.Log.Error("failed to create url service", zap.String("error", err.Error()))
		return nil, err
	}
	return &URLHandler{s, config.ServerConfig.Base.String()}, nil
}

type URLHandler struct {
	service service.Service
	address string
}

func (h URLHandler) createResponseAddress(shortURL string) string {
	return h.address + "/" + shortURL
}

func (h URLHandler) Close() error {
	return h.service.Close()
}

func (h URLHandler) CreateURLWithTextBody(w http.ResponseWriter, r *http.Request) error {
	limitedBody := http.MaxBytesReader(w, r.Body, maxBodySize)
	urlBytes, err := io.ReadAll(limitedBody)
	defer limitedBody.Close()

	if err != nil {
		logger.Log.Error("failed to read request body", zap.String("error", err.Error()))
		return err
	}

	var id string
	id, err = h.service.CreateURL(urlBytes)

	if err != nil {
		logger.Log.Error("failed create url from request body", zap.String("error", err.Error()))
		return err
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.createResponseAddress(id)))

	if err != nil {
		logger.Log.Error("failed to write response body", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (h URLHandler) GetURL(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	url, err := h.service.GetURL(id)

	if err != nil {
		logger.Log.Error("failed to get url by id", zap.String("error", err.Error()))
		return err
	}

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}

func (h URLHandler) CreateShortenWithJSONBody(w http.ResponseWriter, r *http.Request) error {
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("failed to decode body", zap.String("error", err.Error()))
		return err
	}

	if len(req.URL) == 0 {
		logger.Log.Error("unsupported empty body in CreateShorten request")
		return errUnsupportedBody
	}

	id, err := h.service.CreateURL([]byte(req.URL))

	if err != nil {
		logger.Log.Error("failed create url from request body", zap.String("error", err.Error()))
		return err
	}

	resp := models.Response{
		Result: h.createResponseAddress(id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("failed encode body", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (h URLHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	err := h.service.PingDB()

	if err != nil {
		logger.Log.Error("failed to ping db", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h URLHandler) CreateBatch(w http.ResponseWriter, r *http.Request) error {
	var req models.BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("failed to decode body", zap.String("error", err.Error()))
		return err
	}

	if len(req) == 0 {
		logger.Log.Error("unsupported empty body in CreateBatch request")
		return errUnsupportedBody
	}

	var originalURLS []string
	for _, batchRecord := range req {
		originalURLS = append(originalURLS, batchRecord.OriginalURL)
	}

	shortURLS, err := h.service.CreateURLS(originalURLS)
	if err != nil {
		logger.Log.Error("failed create urls", zap.String("error", err.Error()))
		return err
	}

	var response models.BatchResponse
	for i := 0; i < len(req); i++ {
		response = append(response, models.BatchResponseRecord{CorrelationID: req[i].CorrelationID, ShortURL: h.createResponseAddress(shortURLS[i])})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		logger.Log.Error("failed encode body", zap.String("error", err.Error()))
		return err
	}

	return nil
}
