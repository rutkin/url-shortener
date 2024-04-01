package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/middleware"
	"go.uber.org/zap"
)

func NewServer() (*Server, error) {
	handler, err := handlers.NewURLHandler()
	if err != nil {
		logger.Log.Error("failed to create url handler", zap.String("error", err.Error()))
		return nil, err
	}
	return &Server{handler}, nil
}

type Server struct {
	urlHandler *handlers.URLHandler
}

func (s Server) Start() error {
	logger.Log.Info("Running server", zap.String("address", config.ServerConfig.Server.String()))

	err := http.ListenAndServe(config.ServerConfig.Server.String(), s.newRootRouter())

	logger.Log.Info("Server stopped")

	return err
}

func (s Server) Close() error {
	return s.urlHandler.Close()
}

func (s Server) newRootRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.WithCompress)
	userIDRouter := r.With(middleware.WithUserID)
	userIDRouter.Post("/", handlers.NewHandler(s.urlHandler.CreateURLWithTextBody))
	userIDRouter.Get("/{id}", handlers.NewHandler(s.urlHandler.GetURL))
	userIDRouter.Post("/api/shorten", handlers.NewHandler(s.urlHandler.CreateShortenWithJSONBody))
	userIDRouter.Post("/api/shorten/batch", handlers.NewHandler(s.urlHandler.CreateBatch))
	userIDRouter.Get("/ping", s.urlHandler.PingDB)
	userIDRouter.Delete("/api/user/urls", handlers.NewHandler(s.urlHandler.DeleteURLS))
	r.With(middleware.WithAuth).Get("/api/user/urls", handlers.NewHandler(s.urlHandler.GetURLS))
	return r
}
