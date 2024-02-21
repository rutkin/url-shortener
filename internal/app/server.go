package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
)

func newRootRouter() http.Handler {
	urlHandler := handlers.NewURLHandler()
	r := chi.NewRouter()
	r.Post("/", handlers.NewHandler(urlHandler.CreateURL))
	r.Get("/{id}", handlers.NewHandler(urlHandler.GetURL))
	return r
}

func StartServer() {
	err := config.ParseFlags()

	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(config.ServerConfig.Server.String(), newRootRouter())

	if err != nil {
		panic(err)
	}
}
