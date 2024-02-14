package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/repository"
	"github.com/rutkin/url-shortener/internal/app/service"
)

func NewURLHandlerRouter() http.Handler {
	repository := repository.NewInMemoryRepository()
	urlService := service.NewURLService(repository)
	urlHandler := urlHandler{urlService, config.ServerConfig.Base.String()}

	r := chi.NewRouter()
	r.Post("/", MakeHandler(urlHandler.CreateURL))
	r.Get("/{id}", MakeHandler(urlHandler.GetURL))

	return r
}

type urlHandler struct {
	service service.Service
	address string
}

func (h urlHandler) CreateURL(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		return errors.New("unsupported Content-Type header, only text/plain; charset=utf-8 allowed")
	}

	if r.URL.Path != "/" {
		return errors.New("unsupported URL path")
	}

	r.Body = http.MaxBytesReader(w, r.Body, 2000)

	urlBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("can not read body")
	}

	var id string
	id, err = h.service.CreateURL(urlBytes)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.address + "/" + id))
	if err != nil {
		return err
	}
	return nil
}

func (h urlHandler) GetURL(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	url, err := h.service.GetURL(id)
	if err != nil {
		return err
	}
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}
