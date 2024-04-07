package repository

import (
	"database/sql"
	"errors"

	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/models"
)

var ErrConflict = errors.New("repository conflict")
var ErrURLDeleted = errors.New("url deleted")

type URLRecord struct {
	ID  string
	URL string
}

type Repository interface {
	CreateURLS(urls []URLRecord, userID string) error
	CreateURL(id string, url string, userID string) error
	GetURL(id string) (string, error)
	GetURLS(userID string) ([]models.URLRecord, error)
	DeleteURLS(urls []string, userID string) error
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
