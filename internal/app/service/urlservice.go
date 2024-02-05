package service

import (
	"fmt"
	"hash/crc32"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/repository"
)

func NewURLService(repository repository.Repository) Service {
	return urlService{repository}
}

type urlService struct {
	repository repository.Repository
}

func (s urlService) CreateURL(urlBytes []byte) (string, error) {
	urlString := string(urlBytes)

	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return "", err
	}

	id := fmt.Sprintf("%X", crc32.ChecksumIEEE(urlBytes))
	err = s.repository.CreateURL(id, urlString)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s urlService) GetURL(id string) (string, error) {
	return s.repository.GetURL(id)
}
