package repository

type Repository interface {
	CreateURL(id string, url string) error
	GetURL(id string) (string, error)
}
