package main

import (
	"github.com/rutkin/url-shortener/internal/app"
)

func main() {
	server, err := app.NewServer()
	if err != nil {
		panic(err)
	}
	defer server.Close()
	server.Start()
}
