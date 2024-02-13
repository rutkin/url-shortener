package app

import (
	"net/http"
	"net/url"

	"github.com/rutkin/url-shortener/internal/app/handlers"
)

var Address = url.URL{Scheme: "http", Host: "localhost:8080"}

func StartServer() {
	r := handlers.NewURLHandlerRouter()
	err := http.ListenAndServe(Address.Host, r)
	if err != nil {
		panic(err)
	}
}
