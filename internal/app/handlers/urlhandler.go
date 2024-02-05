package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/service"
)

type URLHandler interface {
	Main(w http.ResponseWriter, r *http.Request) error
}

func NewURLHandler(service service.Service, address url.URL) URLHandler {
	return urlHandler{service, address}
}

type urlHandler struct {
	service service.Service
	address url.URL
}

func (h urlHandler) createURL(w http.ResponseWriter, r *http.Request) error {
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
	w.Write([]byte(h.address.JoinPath(id).String()))
	return nil
}

func (h urlHandler) getURL(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[1:]

	url, err := h.service.GetURL(id)
	if err != nil {
		return err
	}
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (h urlHandler) Main(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		return h.createURL(w, r)
	}

	if r.Method == http.MethodGet {
		return h.getURL(w, r)
	}

	return errors.New("unsupported mmethod")
}
