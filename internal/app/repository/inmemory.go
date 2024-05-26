package repository

import (
	"errors"
	"sync"

	"github.com/rutkin/url-shortener/internal/app/models"
)

var errURLNotFound = errors.New("URL not found")
var errNotImplemented = errors.New("not implemented")

// create new instance of repository in memory
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

// store urls in memory
func (r *inMemoryRepository) CreateURLS(urlRecords []URLRecord) error {
	r.mu.Lock()
	for _, record := range urlRecords {
		r.urls[record.ID] = urlValue{longURL: record.URL, userID: record.UserID}
	}
	r.mu.Unlock()
	return nil
}

// store url in memory
func (r *inMemoryRepository) CreateURL(urlRecord URLRecord) error {
	r.mu.Lock()
	r.urls[urlRecord.ID] = urlValue{longURL: urlRecord.URL, userID: urlRecord.UserID}
	r.mu.Unlock()

	return nil
}

// get url from memoty
func (r *inMemoryRepository) GetURL(id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	url, ok := r.urls[id]
	if !ok {
		return "", errURLNotFound
	}

	return url.longURL, nil
}

// get urls from memory
func (r *inMemoryRepository) GetURLS(userID string) ([]models.URLRecord, error) {
	return nil, nil
}

// delete urls from memory
func (r *inMemoryRepository) DeleteURLS(urls []string, userID string) error {
	return errNotImplemented
}

// close
func (r *inMemoryRepository) Close() error {
	return nil
}
