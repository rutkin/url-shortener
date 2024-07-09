package service

import (
	"database/sql"
	"errors"
	"fmt"
	"hash/crc32"
	"net/url"
	"sync"

	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/models"
	"github.com/rutkin/url-shortener/internal/app/repository"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// create new instance of url service
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
	return &urlService{db, r, sync.WaitGroup{}}, nil
}

type urlService struct {
	db         *sql.DB
	repository repository.Repository
	wg         sync.WaitGroup
}

func (s *urlService) createShortURL(url []byte) string {
	return fmt.Sprintf("%X", crc32.ChecksumIEEE(url))
}

func (s *urlService) deleteURLSAsync(urls []string, userID string) {
	defer s.wg.Done()
	s.repository.DeleteURLS(urls, userID)
}

// create urls
func (s *urlService) CreateURLS(urls []string, userID string) ([]string, error) {
	var repositoryURLS []repository.URLRecord
	var shortURLS []string
	for _, url := range urls {
		shortURL := s.createShortURL([]byte(url))
		shortURLS = append(shortURLS, shortURL)
		repositoryURLS = append(repositoryURLS, repository.URLRecord{ID: shortURL, URL: url, UserID: userID})
	}

	err := s.repository.CreateURLS(repositoryURLS)
	if err != nil {
		logger.Log.Error("failed to create urls", zap.String("error", err.Error()))
		return nil, err
	}
	return shortURLS, nil
}

// create url
func (s *urlService) CreateURL(urlBytes []byte, userID string) (string, error) {
	urlString := string(urlBytes)

	_, err := url.ParseRequestURI(urlString)

	if err != nil {
		logger.Log.Error("failed to parse url",
			zap.String("url", urlString),
			zap.String("error", err.Error()))
		return "", err
	}

	id := fmt.Sprintf("%X", crc32.ChecksumIEEE(urlBytes))
	err = s.repository.CreateURL(repository.URLRecord{ID: id, URL: urlString, UserID: userID})

	if errors.Is(err, repository.ErrConflict) {
		return id, err
	}
	if err != nil {
		logger.Log.Error("failed to create url",
			zap.String("url", urlString),
			zap.String("error", err.Error()))
		return "", err
	}

	return id, nil
}

// get url
func (s *urlService) GetURL(id string) (string, error) {
	return s.repository.GetURL(id)
}

// get urls
func (s *urlService) GetURLS(userID string) ([]models.URLRecord, error) {
	return s.repository.GetURLS(userID)
}

// get stats
func (s *urlService) GetStats() (models.StatRecord, error) {
	return s.repository.GetStats()
}

// delete urls
func (s *urlService) DeleteURLS(urls []string, userID string) error {
	s.wg.Add(1)
	go s.deleteURLSAsync(urls, userID)
	return nil
}

// ping database
func (s *urlService) PingDB() error {
	return s.db.Ping()
}

// close instance
func (s *urlService) Close() error {
	s.wg.Wait()
	if s.db != nil {
		s.db.Close()
	}
	return s.repository.Close()
}
