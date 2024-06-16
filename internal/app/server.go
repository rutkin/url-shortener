package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// create new instance of server
func NewServer() (*Server, error) {
	handler, err := handlers.NewURLHandler()
	if err != nil {
		logger.Log.Error("failed to create url handler", zap.String("error", err.Error()))
		return nil, err
	}
	return &Server{handler}, nil
}

// server type
type Server struct {
	urlHandler *handlers.URLHandler
}

// start
func (s Server) Start() error {
	logger.Log.Info("Running server", zap.String("address", config.ServerConfig.Server.String()))

	var srv *http.Server

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("HTTP server Shutdown: ", zap.String("error", err.Error()))
		}
		close(idleConnsClosed)
	}()

	var err error
	if config.ServerConfig.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:  autocert.DirCache("cache-dir"),
			Prompt: autocert.AcceptTOS,
		}

		srv = &http.Server{
			Addr:      ":443",
			Handler:   s.newRootRouter(),
			TLSConfig: manager.TLSConfig(),
		}
		err = srv.ListenAndServeTLS("", "")
	} else {
		srv = &http.Server{Addr: config.ServerConfig.Server.String(), Handler: s.newRootRouter()}
		err = srv.ListenAndServe()
	}

	<-idleConnsClosed

	logger.Log.Info("Server stopped")

	return err
}

// close
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
