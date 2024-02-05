package repository

import (
	"errors"
)

func NewInMemoryRepository() Repository {
	return inMemoryRepository{make(map[string]string)}
}

type inMemoryRepository struct {
	urls map[string]string
}

func (r inMemoryRepository) CreateURL(id string, url string) error {
	r.urls[id] = url
	return nil
}

func (r inMemoryRepository) GetURL(id string) (string, error) {
	url, ok := r.urls[id]
	if !ok {
		return "", errors.New("URL not found")
	}
	return url, nil
}
