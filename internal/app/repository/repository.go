package repository

import (
	"database/sql"
	"errors"

	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/models"
)

// error conflict in repository
var ErrConflict = errors.New("repository conflict")

// error url deleted
var ErrURLDeleted = errors.New("url deleted")

// URLRecord - record to store info about URL in repository
type URLRecord struct {
	// ID - short url id
	ID string
	// URL - original URL
	URL string
	// UserID - user id
	UserID string
}

// Repository - interface for store records
type Repository interface {
	// Create
	CreateURLS(urls []URLRecord) error
	CreateURL(urlRecord URLRecord) error
	GetURL(id string) (string, error)
	GetURLS(userID string) ([]models.URLRecord, error)
	DeleteURLS(urls []string, userID string) error
	GetStats() (models.StatRecord, error)
	Close() error
}

// create new instance of repository in config settings
func NewRepository(db *sql.DB) (Repository, error) {
	if db != nil {
		return NewInDatabaseRepository(db)
	}

	if config.ServerConfig.FileStoragePath == "" {
		return NewInMemoryRepository(), nil
	}

	return NewInFileRepository(config.ServerConfig.FileStoragePath)
}
