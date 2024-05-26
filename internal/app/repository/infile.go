package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

type urlRecord struct {
	ShortURL string `json:"shortURL"`
	LongURL  string `json:"longURL"`
	UserID   string `json:"userID"`
}

func NewInFileRepository(filename string) (*inFileRepository, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Failed to open file repository",
			zap.String("filename", filename),
			zap.String("error", err.Error()))
		return nil, err
	}

	urls := make(map[string]urlValue)
	decoder := json.NewDecoder(f)
	for {
		var urlRecord urlRecord
		err = decoder.Decode(&urlRecord)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			logger.Log.Error("Failed to decode url record", zap.String("error", err.Error()))
		}

		urls[urlRecord.ShortURL] = urlValue{longURL: urlRecord.LongURL, userID: urlRecord.UserID}
	}

	return &inFileRepository{inMemoryRepository: &inMemoryRepository{urls: urls}, file: f, encoder: json.NewEncoder(f)}, nil
}

type inFileRepository struct {
	*inMemoryRepository
	file    *os.File
	encoder *json.Encoder
}

func (r *inFileRepository) CreateURLS(urls []URLRecord) error {
	err := r.inMemoryRepository.CreateURLS(urls)
	if err != nil {
		return err
	}
	r.encoder.Encode(urls)
	return nil
}

func (r *inFileRepository) CreateURL(urlRecord URLRecord) error {
	err := r.inMemoryRepository.CreateURL(urlRecord)
	if err != nil {
		return err
	}

	return r.encoder.Encode(urlRecord)
}

func (r *inFileRepository) Close() error {
	return r.file.Close()
}
