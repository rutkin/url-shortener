package service

type Service interface {
	CreateURL(url []byte) (string, error)
	GetURL(id string) (string, error)
}
