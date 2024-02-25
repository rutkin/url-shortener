package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"go.uber.org/zap"
)

func newRootRouter() http.Handler {
	urlHandler := handlers.NewURLHandler()
	r := chi.NewRouter()
	r.Post("/", handlers.NewHandler(urlHandler.CreateURL))
	r.Get("/{id}", handlers.NewHandler(urlHandler.GetURL))
	r.Post("/api/shorten", handlers.NewHandler(urlHandler.CreateShorten))
	return r
}

func StartServer() {
	err := config.ParseFlags()

	if err != nil {
		panic(err)
	}

	err = logger.Initialize(config.ServerConfig.LogLevel)

	if err != nil {
		panic(err)
	}

	logger.Log.Info("Running server", zap.String("address", config.ServerConfig.Server.String()))

	err = http.ListenAndServe(config.ServerConfig.Server.String(), logger.WithLogging(newRootRouter()))

	if err != nil {
		panic(err)
	}
}
