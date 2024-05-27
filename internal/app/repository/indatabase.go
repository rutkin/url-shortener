package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/models"
	"go.uber.org/zap"
)

// create new instance of database repository
func NewInDatabaseRepository(db *sql.DB) (*inDatabaseRepository, error) {
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Error("Failed to create transaction", zap.String("error", err.Error()))
		return nil, err
	}

	defer tx.Rollback()

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS shortener (shortURL VARCHAR (50) UNIQUE NOT NULL, LongURL VARCHAR (1000) NOT NULL, userID VARCHAR (50) NOT NULL, deleted BOOLEAN NOT NULL)")
	if err != nil {
		logger.Log.Error("Failed to create table", zap.String("error", err.Error()))
		return nil, err
	}

	tx.Exec(`CREATE INDEX IF NOT EXISTS long_url_idx ON shortener (LongURL)`)
	err = tx.Commit()
	if err != nil {
		logger.Log.Error("Failed to prepare db", zap.String("error", err.Error()))
		tx.Rollback()
		return nil, err
	}
	return &inDatabaseRepository{db}, nil
}

type inDatabaseRepository struct {
	db *sql.DB
}

// store urls in db
func (r *inDatabaseRepository) CreateURLS(urls []URLRecord) error {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error("Failed to create transaction", zap.String("error", err.Error()))
		return err
	}

	for _, url := range urls {
		_, err = tx.Exec("INSERT INTO shortener (shortURL, LongURL, userID, deleted) Values ($1, $2, $3, FALSE);", url.ID, url.URL, url.UserID)
		if err != nil {
			logger.Log.Error("Failed to create url", zap.String("error", err.Error()))
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// store url in db
func (r *inDatabaseRepository) CreateURL(urlRecord URLRecord) error {
	_, err := r.db.Exec("INSERT INTO shortener (shortURL, LongURL, userID, deleted) Values ($1, $2, $3, FALSE)", urlRecord.ID, urlRecord.URL, urlRecord.UserID)

	if err != nil {
		logger.Log.Error("Failed to insert in table", zap.String("error", err.Error()))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}

	return nil
}

// get url from db
func (r *inDatabaseRepository) GetURL(id string) (string, error) {
	row := r.db.QueryRow("SELECT LongURL, deleted FROM shortener WHERE shortURL=$1;", id)
	var longURL string
	var deleted bool
	err := row.Scan(&longURL, &deleted)
	if err != nil {
		logger.Log.Error("Failed to select", zap.String("error", err.Error()))
		return "", err
	}
	if deleted {
		return "", ErrURLDeleted
	}
	return longURL, nil
}

// get urls from db
func (r *inDatabaseRepository) GetURLS(userID string) ([]models.URLRecord, error) {
	rows, err := r.db.Query("SELECT shortURL, LongURL FROM shortener WHERE userID=$1;", userID)
	if err != nil {
		logger.Log.Error("Failed to get urls from db", zap.String("error", err.Error()))
		return nil, err
	}

	var result []models.URLRecord

	for rows.Next() {
		err := rows.Err()
		if err != nil {
			logger.Log.Error("Failed to iterate db", zap.String("error", err.Error()))
			return nil, err
		}
		var urlRecord models.URLRecord
		if err := rows.Scan(&urlRecord.ShortURL, &urlRecord.OriginalURL); err != nil {
			logger.Log.Error("Failed to scan get urls result", zap.String("error", err.Error()))
			return nil, err
		}
		result = append(result, urlRecord)
	}

	return result, nil
}

// delete urls from db
func (r *inDatabaseRepository) DeleteURLS(urls []string, userID string) error {
	query := `
		UPDATE shortener SET deleted = TRUE
		WHERE userID=$1 AND shortURL = ANY($2);`
	_, err := r.db.Exec(query, userID, pq.Array(&urls))
	if err != nil {
		logger.Log.Error("Failed to delete urls from db", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// close db
func (r *inDatabaseRepository) Close() error {
	return nil
}
