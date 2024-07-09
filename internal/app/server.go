package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
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
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		s.startGRPC()
		wg.Done()
	}()
	go func() {
		s.startHTTP()
		wg.Done()
	}()
	wg.Wait()
	return nil
}

// start http server
func (s Server) startHTTP() error {
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

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// start grpc server
func (s Server) startGRPC() error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.Log.Error("failed to listen tcp server", zap.String("error", err.Error()))
		return err
	}
	grpcServer := grpc.NewServer()
	grpcHandler, err := handlers.NewGRPCHandler()
	if err != nil {
		return err
	}
	handlers.RegisterGRPCHandlerServer(grpcServer, grpcHandler)
	if err := grpcServer.Serve(listen); err != nil {
		logger.Log.Error("failed to serve grpc", zap.String("error", err.Error()))
		return err
	}
	return nil
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
	userIDRouter.Get("/api/internal/stats", handlers.NewHandler(s.urlHandler.GetStats))
	r.With(middleware.WithAuth).Get("/api/user/urls", handlers.NewHandler(s.urlHandler.GetURLS))
	return r
}
