package app

import (
	"net/http"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/handlers"
	"github.com/rutkin/url-shortener/internal/app/repository"
	"github.com/rutkin/url-shortener/internal/app/service"
)

type Server struct {
	urlHandler handlers.URLHandler
}

var Address = url.URL{Scheme: "http", Host: "localhost:8080"}

func StartServer() {
	repository := repository.NewInMemoryRepository()
	urlService := service.NewURLService(repository)
	urlHandler := handlers.NewURLHandler(urlService, Address)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.MakeHandler(urlHandler.Main))

	err := http.ListenAndServe(Address.Host, mux)
	if err != nil {
		panic(err)
	}
}
