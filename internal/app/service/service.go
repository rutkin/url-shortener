package service

type Service interface {
	CreateURLS(urls []string) ([]string, error)
	CreateURL(url []byte) (string, error)
	GetURL(id string) (string, error)
	PingDB() error
	Close() error
}
