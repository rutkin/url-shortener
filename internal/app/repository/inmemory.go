package repository

import (
	"errors"
	"sync"
)

var errURLNotFound = errors.New("URL not found")

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

func (r *inMemoryRepository) GetURL(id string, userID string) (string, error) {
	r.mu.RLock()
	url, ok := r.urls[userID][id]
	r.mu.RUnlock()

	if !ok {
		return "", errURLNotFound
	}

	return url, nil
}

func (r *inMemoryRepository) GetURLS(userID string) ([]string, error) {
	return nil, nil
}

func (r *inMemoryRepository) Close() error {
	return nil
}
