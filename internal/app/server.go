package app

import (
	"net/http"

	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
)

func StartServer() {
	err := config.ServerConfig.ParseFlags()

	if err != nil {
		panic(err)
	}

	r := handlers.NewURLHandlerRouter()
	err = http.ListenAndServe(config.ServerConfig.Server.String(), r)

	if err != nil {
		panic(err)
	}
}
