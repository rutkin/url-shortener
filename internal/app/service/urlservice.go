package service

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/repository"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewURLService() (*urlService, error) {
	var db *sql.DB
	if len(config.ServerConfig.DatabaseDSN) > 0 {
		newDB, err := sql.Open("pgx", config.ServerConfig.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		db = newDB
	}

	r, err := repository.NewRepository(db)
	if err != nil {
		return nil, err
	}
	return &urlService{db, r}, nil
}

type urlService struct {
	db         *sql.DB
	repository repository.Repository
}

func (s *urlService) createShortURL(url []byte) string {
	return fmt.Sprintf("%X", crc32.ChecksumIEEE(url))
}

func (s *urlService) CreateURLS(urls []string) ([]string, error) {
	var repositoryURLS []repository.URLRecord
	var shortURLS []string
	for _, url := range urls {
		shortURL := s.createShortURL([]byte(url))
		shortURLS = append(shortURLS, shortURL)
		repositoryURLS = append(repositoryURLS, repository.URLRecord{ID: shortURL, URL: url})
	}

	err := s.repository.CreateURLS(repositoryURLS)
	if err != nil {
		logger.Log.Error("failed to create urls", zap.String("error", err.Error()))
		return nil, err
	}
	return shortURLS, nil
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

func (s *urlService) PingDB() error {
	return s.db.Ping()
}

func (s *urlService) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return s.repository.Close()
}
