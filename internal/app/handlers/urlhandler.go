package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/models"
	"github.com/rutkin/url-shortener/internal/app/repository"
	"github.com/rutkin/url-shortener/internal/app/service"
	"go.uber.org/zap"
)

var errUnsupportedBody = errors.New("unsupported body")
var errInvalidContext = errors.New("invalid context")
var errAccessDenied = errors.New("access denied")
var maxBodySize = int64(2000)

// create new instance of url handler
func NewURLHandler() (*URLHandler, error) {
	s, err := service.NewURLService()
	if err != nil {
		logger.Log.Error("failed to create url service", zap.String("error", err.Error()))
		return nil, err
	}
	return &URLHandler{s, config.ServerConfig.Base.String()}, nil
}

// url handler type
type URLHandler struct {
	service service.Service
	address string
}

func (h URLHandler) createResponseAddress(shortURL string) string {
	return h.address + "/" + shortURL
}

func (h URLHandler) writeURLBodyInText(w http.ResponseWriter, shortURL string, statusCode int) error {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	_, err := w.Write([]byte(h.createResponseAddress(shortURL)))
	if err != nil {
		logger.Log.Error("failed to write response body", zap.String("error", err.Error()))
	}
	return err
}

func (h URLHandler) writeURLBodyInJSON(w http.ResponseWriter, shortURL string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := models.Response{
		Result: h.createResponseAddress(shortURL),
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Error("failed encode body", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (h URLHandler) getUserID(context context.Context) (string, error) {
	userID := context.Value(service.UserIDKey)
	if userID == nil {
		logger.Log.Error("userID value does not exists in context")
		return "", errInvalidContext
	}

	return userID.(string), nil
}

// must to call after create, to avoid memory leak
func (h URLHandler) Close() error {
	return h.service.Close()
}

// create short url with text body
func (h URLHandler) CreateURLWithTextBody(w http.ResponseWriter, r *http.Request) error {
	limitedBody := http.MaxBytesReader(w, r.Body, maxBodySize)
	urlBytes, err := io.ReadAll(limitedBody)
	defer limitedBody.Close()

	if err != nil {
		logger.Log.Error("failed to read request body", zap.String("error", err.Error()))
		return err
	}

	userID, err := h.getUserID(r.Context())
	if err != nil {
		return err
	}

	var id string
	id, err = h.service.CreateURL(urlBytes, userID)

	if errors.Is(err, repository.ErrConflict) {
		writeErr := h.writeURLBodyInText(w, id, http.StatusConflict)
		if writeErr != nil {
			return writeErr
		}
		return err
	}

	if err != nil {
		logger.Log.Error("failed create url from request body", zap.String("error", err.Error()))
		return err
	}

	return h.writeURLBodyInText(w, id, http.StatusCreated)
}

// get url by short id
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

// delete batch of urls
func (h URLHandler) DeleteURLS(w http.ResponseWriter, r *http.Request) error {
	var urls []string
	if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
		logger.Log.Error("failed to decode body", zap.String("error", err.Error()))
		return err
	}

	userID, err := h.getUserID(r.Context())
	if err != nil {
		return err
	}

	err = h.service.DeleteURLS(urls, userID)
	if err != nil {
		logger.Log.Error("failed to delete urls", zap.String("error", err.Error()))
		return errAccessDenied
	}

	w.WriteHeader(http.StatusAccepted)

	return nil
}

// get batch of urls
func (h URLHandler) GetURLS(w http.ResponseWriter, r *http.Request) error {
	userID, err := h.getUserID(r.Context())
	if err != nil {
		return err
	}

	urls, err := h.service.GetURLS(userID)
	for k, v := range urls {
		urls[k].ShortURL = h.createResponseAddress(v.ShortURL)
	}

	if err != nil {
		logger.Log.Error("failed to get urls by user id", zap.String("error", err.Error()))
		return err
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(urls); err != nil {
		logger.Log.Error("failed encode body", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// create short url with json body
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

	userID, err := h.getUserID(r.Context())
	if err != nil {
		return err
	}

	id, err := h.service.CreateURL([]byte(req.URL), userID)

	if errors.Is(err, repository.ErrConflict) {
		writeErr := h.writeURLBodyInJSON(w, id, http.StatusConflict)
		if writeErr != nil {
			return writeErr
		}
		return err
	}

	if err != nil {
		logger.Log.Error("failed create url from request body", zap.String("error", err.Error()))
		return err
	}

	return h.writeURLBodyInJSON(w, id, http.StatusCreated)
}

// pind database
func (h URLHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	err := h.service.PingDB()

	if err != nil {
		logger.Log.Error("failed to ping db", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// create batch of short url
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

	userID, err := h.getUserID(r.Context())

	if err != nil {
		return err
	}

	shortURLS, err := h.service.CreateURLS(originalURLS, userID)

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
