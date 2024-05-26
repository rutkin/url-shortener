package service

import "github.com/rutkin/url-shortener/internal/app/models"

type contextKey string

// key used for set/get user id from context
const UserIDKey contextKey = "userID"

// service interface that implement logic
type Service interface {
	CreateURLS(urls []string, userID string) ([]string, error)
	CreateURL(url []byte, userID string) (string, error)
	GetURL(id string) (string, error)
	GetURLS(userID string) ([]models.URLRecord, error)
	DeleteURLS(urls []string, userID string) error
	PingDB() error
	Close() error
}
