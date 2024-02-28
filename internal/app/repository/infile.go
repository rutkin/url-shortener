package repository

import (
	"encoding/json"
	"os"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

func NewInFileRepository(filename string) (*inFileRepository, error) {
	urls := make(map[string]string)

	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Failed to open file repository",
			zap.String("filename", filename),
			zap.String("error", err.Error()))
		return nil, err
	}
	defer f.Close()

	st, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	if st.Size() == 0 {
		return &inFileRepository{inMemoryRepository: NewInMemoryRepository(), filename: filename}, nil
	}

	err = json.NewDecoder(f).Decode(&urls)
	if err != nil {
		logger.Log.Error("Failed to decode file to urls",
			zap.String("filename", filename),
			zap.String("error", err.Error()))
		return nil, err
	}
	return &inFileRepository{inMemoryRepository: &inMemoryRepository{urls: urls}, filename: filename}, nil
}

type inFileRepository struct {
	*inMemoryRepository
	filename string
}

func (r *inFileRepository) Close() error {
	f, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Failed to open file repository",
			zap.String("filename", r.filename),
			zap.String("error", err.Error()))
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(r.urls)
}
