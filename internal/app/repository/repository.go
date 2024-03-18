package repository

import (
	"database/sql"

	"github.com/rutkin/url-shortener/internal/app/config"
)

type URLRecord struct {
	ID  string
	URL string
}

type Repository interface {
	CreateURLS(urls []URLRecord) error
	CreateURL(id string, url string) error
	GetURL(id string) (string, error)
	Close() error
}

func NewRepository(db *sql.DB) (Repository, error) {
	if db != nil {
		return NewInDatabaseRepository(db)
	}

	if config.ServerConfig.FileStoragePath == "" {
		return NewInMemoryRepository(), nil
	}

	return NewInFileRepository(config.ServerConfig.FileStoragePath)
}
