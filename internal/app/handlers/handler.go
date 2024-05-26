package handlers

import (
	"errors"
	"net/http"

	"github.com/rutkin/url-shortener/internal/app/repository"
)

// wrapper function convert error to http error status
func NewHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if errors.Is(err, repository.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, repository.ErrURLDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
