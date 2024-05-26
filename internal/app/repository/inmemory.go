package repository

import (
	"errors"
	"sync"

	"github.com/rutkin/url-shortener/internal/app/models"
)

var errURLNotFound = errors.New("URL not found")
var errNotImplemented = errors.New("not implemented")

func NewInMemoryRepository() *inMemoryRepository {
	res := new(inMemoryRepository)
	res.urls = make(map[string]urlValue)

	return res
}

type urlValue struct {
	longURL string
	userID  string
}

type inMemoryRepository struct {
	urls map[string]urlValue // [shortURL, (longURL, userID)]
	mu   sync.RWMutex
}

func (r *inMemoryRepository) CreateURLS(urlRecords []URLRecord) error {
	r.mu.Lock()
	for _, record := range urlRecords {
		r.urls[record.ID] = urlValue{longURL: record.URL, userID: record.UserID}
	}
	r.mu.Unlock()
	return nil
}

func (r *inMemoryRepository) CreateURL(urlRecord URLRecord) error {
	r.mu.Lock()
	r.urls[urlRecord.ID] = urlValue{longURL: urlRecord.URL, userID: urlRecord.UserID}
	r.mu.Unlock()

	return nil
}

func (r *inMemoryRepository) GetURL(id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	url, ok := r.urls[id]
	if !ok {
		return "", errURLNotFound
	}

	return url.longURL, nil
}

func (r *inMemoryRepository) GetURLS(userID string) ([]models.URLRecord, error) {
	return nil, nil
}

func (r *inMemoryRepository) DeleteURLS(urls []string, userID string) error {
	return errNotImplemented
}

func (r *inMemoryRepository) Close() error {
	return nil
}
