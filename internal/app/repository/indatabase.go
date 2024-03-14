package repository

import (
	"database/sql"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

func NewInDatabaseRepository(db *sql.DB) (*inDatabaseRepository, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXIST shortener (shortURL VARCHAR (50) UNIQUE NOT NULL, LongURL VARCHAR (1000) NOT NULL);")
	if err != nil {
		logger.Log.Error("Failed to create table", zap.String("error", err.Error()))
		return nil, err
	}

	return &inDatabaseRepository{db}, nil
}

type inDatabaseRepository struct {
	db *sql.DB
}

func (r *inDatabaseRepository) CreateURL(id string, url string) error {
	_, err := r.db.Exec("INSERT INTO shortener (shortURL, LongURL) Values ($1, $2);", id, url)
	if err != nil {
		logger.Log.Error("Failed to insert in table", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *inDatabaseRepository) GetURL(id string) (string, error) {
	row := r.db.QueryRow("SELECT LongURL FROM shortener WHERE shortURL=$1;", id)
	var longURL string
	err := row.Scan(longURL)
	if err != nil {
		logger.Log.Error("Failed to select", zap.String("error", err.Error()))
		return "", err
	}
	return longURL, nil
}

func (r *inDatabaseRepository) Close() error {
	return nil
}
