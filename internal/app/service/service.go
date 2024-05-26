package service

import "github.com/rutkin/url-shortener/internal/app/models"

type contextKey string

// key used for set/get user id from context
const UserIDKey contextKey = "userID"

// service interface that implement logic
type Service interface {
	// create urls
	CreateURLS(urls []string, userID string) ([]string, error)
	// create url
	CreateURL(url []byte, userID string) (string, error)
	// get url
	GetURL(id string) (string, error)
	// get urls
	GetURLS(userID string) ([]models.URLRecord, error)
	// delete urls
	DeleteURLS(urls []string, userID string) error
	// ping database
	PingDB() error
	// close
	Close() error
}
