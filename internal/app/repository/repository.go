package repository

import (
	"github.com/rutkin/url-shortener/internal/app/config"
)

type Repository interface {
	CreateURL(id string, url string) error
	GetURL(id string) (string, error)
	Close() error
}

func NewRepository() (Repository, error) {
	if config.ServerConfig.FileStoragePath == "" {
		return NewInMemoryRepository(), nil
	}

	return NewInFileRepository(config.ServerConfig.FileStoragePath)
}
