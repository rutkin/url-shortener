package repository

import (
	"database/sql"
	"errors"

	"github.com/rutkin/url-shortener/internal/app/config"
)

var ErrConflict = errors.New("repository conflict")

type URLRecord struct {
	ID  string
	URL string
}

type Repository interface {
	CreateURLS(urls []URLRecord, userID string) error
	CreateURL(id string, url string, userID string) error
	GetURL(id string, userID string) (string, error)
	GetURLS(userID string) ([]string, error)
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
