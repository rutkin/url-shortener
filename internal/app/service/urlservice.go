package service

import (
	"fmt"
	"hash/crc32"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/repository"
	"go.uber.org/zap"
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
		logger.Log.Error("failed to parse url",
			zap.String("url", urlString),
			zap.String("error", err.Error()))
		return "", err
	}

	id := fmt.Sprintf("%X", crc32.ChecksumIEEE(urlBytes))
	err = s.repository.CreateURL(id, urlString)

	if err != nil {
		logger.Log.Error("failed to create url",
			zap.String("url", urlString),
			zap.String("error", err.Error()))
		return "", err
	}

	return id, nil
}

func (s *urlService) GetURL(id string) (string, error) {
	return s.repository.GetURL(id)
}

func (s *urlService) Close() error {
	return s.repository.Close()
}
