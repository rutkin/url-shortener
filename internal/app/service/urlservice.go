package service

import (
	"fmt"
	"hash/crc32"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/repository"
)

func NewURLService() (*urlService, error) {
	r, err := repository.NewRepository()
	if err != nil {
		return nil, err
	}
	return &urlService{r}, nil
}

type urlService struct {
	repository repository.Repository
}

func (s *urlService) CreateURL(urlBytes []byte) (string, error) {
	urlString := string(urlBytes)

	_, err := url.ParseRequestURI(urlString)

	if err != nil {
		return "", fmt.Errorf("failed to parse url '%s': %w", urlString, err)
	}

	id := fmt.Sprintf("%X", crc32.ChecksumIEEE(urlBytes))
	err = s.repository.CreateURL(id, urlString)

	if err != nil {
		return "", fmt.Errorf("failed to create url '%s': %w", urlString, err)
	}

	return id, nil
}

func (s *urlService) GetURL(id string) (string, error) {
	return s.repository.GetURL(id)
}

func (s *urlService) Close() error {
	return s.repository.Close()
}
