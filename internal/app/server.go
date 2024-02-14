package app

import (
	"flag"
	"net/http"

	"github.com/rutkin/url-shortener/internal/app/config"
	"github.com/rutkin/url-shortener/internal/app/handlers"
)

func StartServer() {
	flag.Parse()
	r := handlers.NewURLHandlerRouter()
	err := http.ListenAndServe(config.ServerConfig.String(), r)
	if err != nil {
		panic(err)
	}
}
