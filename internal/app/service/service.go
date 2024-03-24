package service

type contextKey string

const UserIDKey contextKey = "userID"

type Service interface {
	CreateURLS(urls []string, userID string) ([]string, error)
	CreateURL(url []byte, userID string) (string, error)
	GetURL(id string, userID string) (string, error)
	PingDB() error
	Close() error
}
