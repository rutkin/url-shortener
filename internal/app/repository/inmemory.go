package repository

import (
	"errors"
	"sync"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

var errURLNotFound = errors.New("URL not found")

func NewInMemoryRepository() *inMemoryRepository {
	res := new(inMemoryRepository)
	res.urls = make(map[string]string)

	return res
}

type inMemoryRepository struct {
	urls map[string]string
	mu   sync.RWMutex
}

func (r *inMemoryRepository) CreateURLS(urls []URLRecord) error {
	for _, url := range urls {
		err := r.CreateURL(url.ID, url.URL)
		if err != nil {
			logger.Log.Error("Failed to create url", zap.String("error", err.Error()))
			return err
		}
	}
	return nil
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

func (r *inMemoryRepository) Close() error {
	return nil
}
