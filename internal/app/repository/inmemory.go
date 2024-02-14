package repository

import (
	"errors"
	"sync"
)

var errURLNotFound = errors.New("URL not found")

func NewInMemoryRepository() Repository {
	res := new(inMemoryRepository)
	res.urls = make(map[string]string)

	return res
}

type inMemoryRepository struct {
	urls map[string]string
	mu   sync.RWMutex
}

func (r *inMemoryRepository) CreateURL(id string, url string) error {
	r.mu.Lock()
	r.urls[id] = url
	r.mu.Unlock()

	return nil
}

func (r *inMemoryRepository) GetURL(id string) (string, error) {
	r.mu.RLock()
	url, ok := r.urls[id]
	r.mu.RUnlock()

	if !ok {
		return "", errURLNotFound
	}

	return url, nil
}
