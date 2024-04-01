package repository

import (
	"errors"
	"sync"

	"github.com/rutkin/url-shortener/internal/app/models"
)

var errURLNotFound = errors.New("URL not found")
var errNotImplemented = errors.New("Not implemented")

func NewInMemoryRepository() *inMemoryRepository {
	res := new(inMemoryRepository)
	res.urls = make(map[string]map[string]string)

	return res
}

type inMemoryRepository struct {
	urls map[string]map[string]string // [userID, [shortURL, longURL]]
	mu   sync.RWMutex
}

func (r *inMemoryRepository) CreateURLS(urls []URLRecord, userID string) error {
	r.mu.Lock()
	for _, url := range urls {
		if r.urls[userID] == nil {
			r.urls[userID] = make(map[string]string)
		}
		r.urls[userID][url.ID] = url.URL
	}
	r.mu.Unlock()
	return nil
}

func (r *inMemoryRepository) CreateURL(id string, url string, userID string) error {
	r.mu.Lock()
	if r.urls[userID] == nil {
		r.urls[userID] = make(map[string]string)
	}
	r.urls[userID][id] = url
	r.mu.Unlock()

	return nil
}

func (r *inMemoryRepository) GetURL(id string) (string, error) {
	var url string
	r.mu.RLock()
	for _, userURL := range r.urls {
		var ok bool
		if url, ok = userURL[id]; ok {
			break
		}
	}
	r.mu.RUnlock()

	if len(url) == 0 {
		return "", errURLNotFound
	}

	return url, nil
}

func (r *inMemoryRepository) GetURLS(userID string) ([]models.URLRecord, error) {
	return nil, nil
}

func (r *inMemoryRepository) GetURLSUserID(urls []string) ([]string, error) {
	return nil, errNotImplemented
}

func (r *inMemoryRepository) Close() error {
	return nil
}
